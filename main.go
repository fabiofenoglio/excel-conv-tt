package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fabiofenoglio/excelconv/logger"
	"github.com/jessevdk/go-flags"

	"github.com/fabiofenoglio/excelconv/model"
	"github.com/fabiofenoglio/excelconv/services"
)

func main() {
	l := logger.GetLogger()
	err := run()
	if err != nil {
		l.Error(err.Error())
		time.Sleep(time.Second * 10)
		os.Exit(1)
	}
}

func run() error {
	var args model.Args
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
		executionErr = services.Run(args)
		return
	}()

	return err
}
