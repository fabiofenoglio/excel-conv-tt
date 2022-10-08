package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/logger"
	"github.com/fabiofenoglio/excelconv/parser"
	"github.com/fabiofenoglio/excelconv/writer"
	excelwriter "github.com/fabiofenoglio/excelconv/writer/excel"
	"github.com/fabiofenoglio/excelconv/writer/json"
)

func main() {
	l := logger.GetLogger()
	err := runWithPanicProtection()
	if err != nil {
		l.Error(err.Error())
		time.Sleep(time.Second * 10)
		os.Exit(1)
	}
}

func runWithPanicProtection() error {
	var args config.Args
	_, err := flags.Parse(&args)
	if err != nil {
		return err
	}

	err = func() (executionErr error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				executionErr = fmt.Errorf("[panic] %v", recovered)
			}
		}()
		executionErr = run(args)
		return
	}()

	return err
}

func run(args config.Args) error {
	input := args.PositionalArgs.InputFile
	if input == "" {
		return errors.New("missing input file")
	}
	log := logger.GetLogger()
	if args.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}
	if args.StdOut {
		log.SetOutput(os.Stderr)
	}

	// parse the specified input file

	parsed, err := parser.Parse(input, log)
	if err != nil {
		return fmt.Errorf("error reading from input file: %w", err)
	}
	log.Infof("found %d rows", len(parsed.Rows))

	writer, err := pickWriter(args)
	if err != nil {
		return err
	}

	outputFile := writer.ComputeDefaultOutputFile(input)
	bytes, err := writer.Write(parsed, log)

	if err != nil {
		return fmt.Errorf("error running writer: %w", err)
	}

	if args.StdOut {
		f := bufio.NewWriter(os.Stdout)
		_, err := f.Write(bytes)
		if err != nil {
			return fmt.Errorf("error writing to stdout: %w", err)
		}
		_ = f.Flush()

	} else {

		log.Debugf("writing to output file %s", outputFile)
		err = ioutil.WriteFile(outputFile, bytes, 0755)
		if err != nil {
			return fmt.Errorf("error saving to output file %s: %w", outputFile, err)
		}
		log.Infof("saved to output file %s", outputFile)

	}

	return nil
}

func pickWriter(arg config.Args) (writer.Writer, error) {
	switch arg.Format {
	case "excel":
		return &excelwriter.WriterImpl{}, nil
	case "json":
		return &json.WriterImpl{}, nil
	}

	return nil, fmt.Errorf("%s is not a valid output format (no writer available)", arg.Format)
}
