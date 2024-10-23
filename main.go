package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/blang/semver"
	aggregator2 "github.com/fabiofenoglio/excelconv/aggregator/v2"
	parser2 "github.com/fabiofenoglio/excelconv/parser/v2"
	"github.com/fabiofenoglio/excelconv/reader/v2"
	excelwriter2 "github.com/fabiofenoglio/excelconv/writer/excel/v2"
	"github.com/fabiofenoglio/go-github-selfupdate/selfupdate"
	"github.com/getsentry/sentry-go"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/fabiofenoglio/excelconv/config"
	"github.com/fabiofenoglio/excelconv/logger"
	"github.com/fabiofenoglio/excelconv/writer"
	jsonwriter2 "github.com/fabiofenoglio/excelconv/writer/json/v2"
)

func logErrorWaitAndExit(err error) {
	logger.GetLogger().Error(err.Error())
	//time.Sleep(time.Second * 20)
	waitForUserPress()
	os.Exit(1)
}

func waitForUserPress() {
	fmt.Println("\nPREMI UN TASTO PER CONTINUARE ...")
	_, _ = fmt.Scanln()
}

func main() {
	envConfig, err := config.Get()
	if err != nil {
		logErrorWaitAndExit(err)
		return
	}

	runWithEnv(envConfig)
}

func setupSentry(envConfig config.EnvConfig) (bool, error) {
	if envConfig.SentryDSN == "" {
		return false, nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              envConfig.SentryDSN,
		TracesSampleRate: 1.0,
		Environment:      envConfig.Environment,
	})
	if err != nil {
		return false, errors.Wrap(err, "error setting up error watch")
	}

	return true, nil
}

func withCustomSentryScope(task func(scope *sentry.Scope)) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetExtra("githubVersion", githubVersion)
		task(scope)
	})
}

func attemptAutoUpdate(envConfig config.EnvConfig, sentryAvailable bool) {
	log := logger.GetLogger()

	defer func() {
		if recovered := recover(); recovered != nil {
			log.Errorf("unexpected error during autoupdate: %v", recovered)
		}
	}()

	if envConfig.SkipAutoUpdate {
		log.Info("skipping auto update because of env configuration")
		return
	}

	fail := func(e error, context string) {
		if sentryAvailable {
			withCustomSentryScope(func(scope *sentry.Scope) {
				sentry.CaptureException(errors.Wrap(e, context))
			})
			log.Info("l'errore è stato segnalato automaticamente")
		}
		log.WithError(e).Error(context)
	}

	v := semver.MustParse(githubVersion)
	up, err := selfupdate.NewUpdater(selfupdate.Config{
		BinaryName: githubBinaryName,
	})

	if err != nil {
		fail(err, "errore nella preparazione per la verifica della versione")
		return
	}

	remoteVersion, found, err := up.DetectLatest(githubRepo)
	if err != nil {
		fail(err, "errore nella verifica della versione")
		return
	}
	if !found {
		fail(errors.New("nessuna sorgente disponibile per la verifica dell'aggiornamento"),
			"errore nella verifica della versione")
		return
	}

	log.Debugf("la versione remota è %s, quella locale %s", remoteVersion.Version.String(), v.String())

	if remoteVersion.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up-to-date.
		log.Debugf("già aggiornato alla versione più recente: %s", githubVersion)
		return
	}

	log.Infof("ATTENZIONE: verifica dell'aggiornamento in corso ...")

	latest, err := up.UpdateSelf(v, githubRepo)
	if err != nil {
		fail(err, "errore nell'aggiornamento automatico")
		return
	}

	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up-to-date.
		log.Debugf("già aggiornato alla versione più recente: %s", githubVersion)
		return
	}

	log.Info("******************************")
	log.Infof("il programma è stato aggiornato automaticamente alla versione più recente: %s", latest.Version)
	fmt.Println("")
	log.Infof(latest.ReleaseNotes)
	fmt.Println("")
	log.Info("ATTENZIONE: è necessario avviare nuovamente il programma per eseguirlo con la nuova versione")

	if sentryAvailable {
		withCustomSentryScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelInfo)
			scope.SetExtra("oldVersion", v.String())
			scope.SetExtra("newVersion", latest.Version.String())
			sentry.CaptureMessage(fmt.Sprintf("automatically updated from version %s to %s",
				v.String(), latest.Version.String()))
		})
	}

	waitForUserPress()
	os.Exit(0)
}

func runWithEnv(envConfig config.EnvConfig) {
	log := logger.GetLogger()

	sentryAvailable, err := setupSentry(envConfig)
	if err != nil {
		log.Error(errors.Wrap(err, "error setting up error watch"))
	}
	if sentryAvailable {
		defer sentry.Flush(2 * time.Second)
	}

	fail := func(e error) {
		if sentryAvailable {
			withCustomSentryScope(func(scope *sentry.Scope) {
				sentry.CaptureException(e)
			})
			log.Info("l'errore è stato segnalato automaticamente")
		}
		logErrorWaitAndExit(e)
	}

	attemptAutoUpdate(envConfig, sentryAvailable)

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

	// parse the specified input file
	ctx := context.Background()
	span := sentry.StartSpan(ctx, "process",
		sentry.WithTransactionName(fmt.Sprintf("process: %s", args.PositionalArgs.InputFile)))
	defer span.Finish()

	err = runWithPanicProtection(span.Context(), args, envConfig, log)
	if err != nil {
		fail(err)
		return
	} else if sentryAvailable {
		withCustomSentryScope(func(scope *sentry.Scope) {
			scope.SetLevel(sentry.LevelInfo)
			scope.SetExtra("arguments", args)
			sentry.CaptureMessage("successful execution")
		})
	}
}

func runWithPanicProtection(ctx context.Context, args config.Args, envConfig config.EnvConfig, logger *logrus.Logger) error {

	err := func() (executionErr error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				executionErr = errors.Errorf("[panic] %v", recovered)
			}
		}()
		executionErr = run(ctx, args, envConfig, logger)
		return
	}()

	return err
}

func run(ctx context.Context, args config.Args, _ config.EnvConfig, log *logrus.Logger) error {
	input := args.PositionalArgs.InputFile
	if input == "" {
		return errors.New("missing input file")
	}

	workflowContext := config.WorkflowContext{
		Context: ctx,
		Logger:  log.WithContext(ctx),
		Config: config.WorkflowContextConfig{
			EnableMissingOperatorsWarning: args.EnableMissingOperatorsWarning,
		},
	}

	span := sentry.StartSpan(ctx, "read")
	readerOutput, err := reader.Execute(workflowContext.ForContext(span.Context()), reader.Input{
		FilePath: input,
	})
	span.Finish()
	if err != nil {
		return err
	}

	span = sentry.StartSpan(ctx, "parse")
	parserOutput, err := parser2.Execute(workflowContext.ForContext(span.Context()), readerOutput)
	span.Finish()
	if err != nil {
		return err
	}

	span = sentry.StartSpan(ctx, "aggregate")
	aggregatorOutput, err := aggregator2.Execute(workflowContext.ForContext(span.Context()), parserOutput)
	span.Finish()
	if err != nil {
		return err
	}

	writer, err := pickWriter(args)
	if err != nil {
		return err
	}

	span = sentry.StartSpan(ctx, "write")
	outputFile := writer.ComputeDefaultOutputFile(input)
	bytes, err := writer.Write(workflowContext.ForContext(span.Context()), aggregatorOutput, parserOutput.Anagraphics)
	span.Finish()
	if err != nil {
		return errors.Wrap(err, "error running writer")
	}

	err = func() error {
		span = sentry.StartSpan(ctx, "dump")
		defer span.Finish()
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
	}()
	if err != nil {
		return err
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
