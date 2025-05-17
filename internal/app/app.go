package app

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/Calyr3x/QuietGrooveBackend/internal/integrations/telegram"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/sirupsen/logrus"
	"sync"
)

type App struct {
	repo        *Registry
	rest        *Rest
	notifier    *telegram.Adapter
	appCron     *AppCron
	controllers *Controllers
	usecases    *Usecases
}

func New(
	ctx context.Context,
	logger logrus.FieldLogger,
	version string,
	config *configuration.Config,
	creds *configuration.Credentials,
) (*App, error) {
	if config == nil {
		return nil, errorspkg.NewErrConstructorDependencies("App", "Config", "nil")
	}
	if creds == nil {
		return nil, errorspkg.NewErrConstructorDependencies("App", "Credentials", "nil")
	}

	repo, err := NewRepo(ctx, creds)
	if err != nil {
		return nil, err
	}

	tgBot, err := telegram.NewAdapter(&creds.TelegramBot)
	if err != nil {
		return nil, err
	}

	usecases, err := NewUsecases(logger, config, repo, tgBot)
	if err != nil {
		return nil, err
	}

	tgBot.RegisterHandlers(usecases.verification)

	controllers, err := NewControllers(logger, usecases)
	if err != nil {
		return nil, err
	}

	restServer, err := NewRest(
		controllers,
		logger,
		config.WebServer,
		version,
	)
	if err != nil {
		return nil, err
	}

	appCron, err := NewAppCron(logger)
	if err != nil {
		return nil, err
	}

	return &App{
		repo:        repo,
		rest:        restServer,
		appCron:     appCron,
		controllers: controllers,
		usecases:    usecases,
		notifier:    tgBot,
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) error {
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = a.rest.Start(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = a.appCron.Start(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.notifier.Run(ctx)
	}()

	return err
}
