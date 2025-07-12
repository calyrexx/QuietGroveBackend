package app

import (
	"context"
	"github.com/calyrexx/QuietGrooveBackend/internal/configuration"
	"github.com/calyrexx/QuietGrooveBackend/internal/integrations/telegram"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"log/slog"
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
	logger *slog.Logger,
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

	tgBot, err := telegram.NewAdapter(&creds.TelegramBot, logger)
	if err != nil {
		return nil, err
	}

	usecases, err := NewUsecases(logger, config, repo, tgBot)
	if err != nil {
		return nil, err
	}

	tgBot.RegisterHandlers(usecases.verification, usecases.reservations)

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

	appCron.Add(config.AppCron.UpdateReservationsStatuses.Spec, usecases.reservations.UpdateStatuses)
	appCron.Add(config.AppCron.GetForReminder.Spec, usecases.reservations.GetForReminder)

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
