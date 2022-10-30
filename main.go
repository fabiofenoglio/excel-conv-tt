package main

import (
	"bufio"
	"context"
	"os"
	"time"

	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
	"github.com/fabiofenoglio/excelconv/reader/v2"
	excelwriter2 "github.com/fabiofenoglio/excelconv/writer/excel/v2"
	"github.com/getsentry/sentry-go"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/logger"
	"github.com/fabiofenoglio/excelconv/writer"
	jsonwriter2 "github.com/fabiofenoglio/excelconv/writer/json/v2"
)

func main() {
	log := logger.GetLogger()

	fail := func(e error) {
		log.Error(e.Error())
		time.Sleep(time.Second * 10)
		os.Exit(1)
	}

	envConfig, err := config.Get()
	if err != nil {
		fail(err)
		return
	}

	sentryAvailable := false
	if envConfig.SentryDSN != "" {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              envConfig.SentryDSN,
			TracesSampleRate: 1.0,
			Environment:      envConfig.Environment,
		})
		if err != nil {
			log.Error(errors.Wrap(err, "error setting up error watch"))
		} else {
			sentryAvailable = true
		}
	}

	if sentryAvailable {
		defer sentry.Flush(2 * time.Second)

		fail = func(e error) {
			log.Error(e.Error())

			sentry.CaptureException(err)
			log.Info("l'errore è stato segnalato automaticamente")

			time.Sleep(time.Second * 10)
			os.Exit(1)
		}
	}

	var args config.Args
	_, err = flags.Parse(&args)
	if err != nil {
		fail(err)
		return
	}

	if args.Verbose {
		log.SetLevel(logrus.DebugLevel)
	}
	if args.StdOut {
		log.SetOutput(os.Stderr)
	}

	err = runWithPanicProtection(args, envConfig, log)
	if err != nil {
		fail(err)
		return
	} else {
		if sentryAvailable {
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelInfo)
				scope.SetExtra("arguments", args)
				sentry.CaptureMessage("successful execution")
			})
		}
	}
}

func runWithPanicProtection(args config.Args, envConfig config.EnvConfig, logger *logrus.Logger) error {

	err := func() (executionErr error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				executionErr = errors.Errorf("[panic] %v", recovered)
			}
		}()
		executionErr = run(args, envConfig, logger)
		return
	}()

	return err
}

func run(args config.Args, _ config.EnvConfig, log *logrus.Logger) error {
	input := args.PositionalArgs.InputFile
	if input == "" {
		return errors.New("missing input file")
	}

	// parse the specified input file
	ctx := context.Background()

	workflowContext := config.WorkflowContext{
		Context: ctx,
		Logger:  log.WithContext(ctx),
		Config: config.WorkflowContextConfig{
			EnableMissingOperatorsWarning: args.EnableMissingOperatorsWarning,
		},
	}

	readerOutput, err := reader.Execute(workflowContext, reader.Input{
		FilePath: input,
	})
	if err != nil {
		return err
	}

	parserOutput, err := parser2.Execute(workflowContext, readerOutput)
	if err != nil {
		return err
	}

	aggregatorOutput, err := aggregator2.Execute(workflowContext, parserOutput)
	if err != nil {
		return err
	}

	writer, err := pickWriter(args)
	if err != nil {
		return err
	}

	outputFile := writer.ComputeDefaultOutputFile(input)
	bytes, err := writer.Write(workflowContext, aggregatorOutput, parserOutput.Anagraphics)

	if err != nil {
		return errors.Wrap(err, "error running writer")
	}

	if args.StdOut {
		f := bufio.NewWriter(os.Stdout)
		_, err := f.Write(bytes)
		if err != nil {
			return errors.Wrap(err, "error writing to stdout")
		}
		_ = f.Flush()

	} else {

		log.Debugf("writing to output file %s", outputFile)
		err = os.WriteFile(outputFile, bytes, 0755)
		if err != nil {
			return errors.Wrapf(err, "error saving to output file %s", outputFile)
		}
		log.Infof("saved to output file %s", outputFile)

	}

	return nil
}

func pickWriter(arg config.Args) (writer.Writer, error) {
	switch arg.Format {
	case "excel":
		return &excelwriter2.WriterImpl{}, nil
	case "json":
		return &jsonwriter2.WriterImpl{}, nil
	}

	return nil, errors.Errorf("%s is not a valid output format (no writer available)", arg.Format)
}
