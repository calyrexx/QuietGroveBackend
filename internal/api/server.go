package api

import (
	"context"
	"net/http"
	"time"
)

type ServerDependencies struct {
	Handler http.Handler
	Config  ServerConfig
}

type ServerConfig struct {
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ShutdownTimeout   time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}

type Server struct {
	srvHTTP         *http.Server
	shutdownTimeout time.Duration
}

func NewServer(dep ServerDependencies) *Server {
	s := &Server{
		srvHTTP: &http.Server{
			Addr:              ":" + dep.Config.Port,
			Handler:           dep.Handler,
			ReadTimeout:       dep.Config.ReadTimeout,
			ReadHeaderTimeout: dep.Config.ReadHeaderTimeout,
			WriteTimeout:      dep.Config.WriteTimeout,
			IdleTimeout:       dep.Config.IdleTimeout,
			MaxHeaderBytes:    dep.Config.MaxHeaderBytes,
		},
		shutdownTimeout: dep.Config.ShutdownTimeout,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	errChan := make(chan error)

	go func() {
		errChan <- s.srvHTTP.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return s.Stop(ctx)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srvHTTP.Shutdown(ctx)
}
