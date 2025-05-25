package main

import (
	"context"
	"fmt"
	"github.com/calyrexx/QuietGrooveBackend/internal/app"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/zeroslog"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const Version = "v0.3.0"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.New(zeroslog.New(
		zeroslog.WithOutput(os.Stderr),
		zeroslog.WithColors(),
		zeroslog.WithTimeFormat("2006-01-02 15:04:05.000 -07:00"),
	))

	logger.Info(fmt.Sprintf("app version: %s", Version))

	config, err := configuration.NewConfig()
	if err != nil {
		logger.Error("newConfig initialization failed", zeroslog.ErrorKey, err)
		return
	}
	config.Version = Version

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

	logger.Info("App has been started!")
	<-ctx.Done()
	logger.Info("Please wait, services are stopping...")
	wg.Wait()
	logger.Info("App is stopped correctly!")
}
