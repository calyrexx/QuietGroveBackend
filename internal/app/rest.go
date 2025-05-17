package app

import (
	"context"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api/handlers"
	"github.com/Calyr3x/QuietGrooveBackend/internal/api/middleware"
	"github.com/Calyr3x/QuietGrooveBackend/internal/configuration"
	"github.com/sirupsen/logrus"
)

type Rest struct {
	httpserver *api.Server
}

func NewRest(
	controllers *Controllers,
	logger logrus.FieldLogger,
	config *configuration.HttpServer,
	version string,
) (*Rest, error) {

	general, err := handlers.NewGeneral(version)
	if err != nil {
		return nil, err
	}

	panicRecoveryMiddleware, err := middleware.NewPanicRecoveryMiddleware(middleware.PanicRecoveryMiddlewareDependencies{
		Logger: logger,
	})
	if err != nil {
		return nil, err
	}

	reservationsHandler, err := handlers.NewReservations(handlers.ReservationsDependencies{
		Controller: controllers.Reservations,
		Logger:     logger,
	})
	if err != nil {
		return nil, err
	}

	housesHandler, err := handlers.NewHouses(handlers.HousesDependencies{
		Controller: controllers.Houses,
		Logger:     logger,
	})
	if err != nil {
		return nil, err
	}

	extrasHandler, err := handlers.NewExtras(handlers.ExtrasDependencies{
		Controller: controllers.Extras,
		Logger:     logger,
	})
	if err != nil {
		return nil, err
	}

	verificationHandler, err := handlers.NewVerification(handlers.VerificationDependencies{
		Controller: controllers.Verification,
	})
	if err != nil {
		return nil, err
	}

	router := api.NewRouter(api.RouterDependencies{
		Handlers: api.Handlers{
			Reservations: reservationsHandler,
			Houses:       housesHandler,
			Extras:       extrasHandler,
			General:      general,
			Verification: verificationHandler,
		},
		Middlewares: api.Middlewares{
			PanicRecovery: panicRecoveryMiddleware.Middleware,
		},
	})

	server := api.NewServer(api.ServerDependencies{
		Handler: router,
		Config: api.ServerConfig{
			Port:              config.Port,
			ReadTimeout:       config.ReadTimeout,
			WriteTimeout:      config.WriteTimeout,
			ShutdownTimeout:   config.ShutdownTimeout,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			IdleTimeout:       config.IdleTimeout,
			MaxHeaderBytes:    config.MaxHeaderBytes,
		},
	})

	return &Rest{
		httpserver: server,
	}, nil
}

func (a *Rest) Start(ctx context.Context) error {
	return a.httpserver.Start(ctx)
}
