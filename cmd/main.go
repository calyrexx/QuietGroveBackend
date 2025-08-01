package main

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/app"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/zeroslog"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const Version = "v0.6.0"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, err := configuration.NewConfig()
	if err != nil {
		println("newConfig initialization failed. error =", err.Error())
		return
	}
	config.Version = Version

	logger := slog.New(zeroslog.New(
		zeroslog.WithOutput(os.Stderr),
		zeroslog.WithColors(),
		zeroslog.WithMinLevel(config.Logger.Level),
		zeroslog.WithTimeFormat("2006-01-02 15:04:05.000 -07:00"),
	))

	logger.Info("Starting app...", "version", Version)

	creds, err := configuration.NewCredentials()
	if err != nil {
		logger.Error("newCredentials initialization failed", zeroslog.ErrorKey, err)
		return
	}

	application, err := app.New(ctx, logger, Version, config, creds)
	if err != nil {
		logger.Error("application.New initialization failed", zeroslog.ErrorKey, err)
		return
	}

	wg := &sync.WaitGroup{}

	err = application.Start(ctx, wg)
	if err != nil {
		logger.Error("application.Start failed", zeroslog.ErrorKey, err)
		return
	}

	logger.Info("App has been started!", "version", Version)
	<-ctx.Done()
	logger.Info("Please wait, services are stopping...")
	wg.Wait()
	logger.Info("App is stopped correctly!")
}
