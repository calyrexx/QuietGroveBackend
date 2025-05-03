package main

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/app"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

const Version = "v0.0.1"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := logrus.New()
	if logger == nil {
		log.Fatal("New logger failed")
	}

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000 -07:00",
	})

	logger.Infof("Version: %s", Version)

	config, err := configuration.NewConfig()
	if err != nil {
		logger.WithError(err).Fatal("NewConfig initialization failed")
	}
	config.Version = Version

	creds, err := configuration.NewCredentials()
	if err != nil {
		logger.WithError(err).Fatal("NewCredentials initialization failed")
	}

	application, err := app.New(ctx, logger, Version, config, creds)
	if err != nil {
		logger.WithError(err).Fatal("application.New initialization failed")
	}

	wg := &sync.WaitGroup{}

	err = application.Start(ctx, wg)
	if err != nil {
		logger.WithError(err).Fatal("application.Start failed")
	}

	logger.Info("App has been started!")
	<-ctx.Done()
	logger.Info("Please wait, services are stopping...")
	wg.Wait()
	logger.Info("App is stopped correctly!")
}
