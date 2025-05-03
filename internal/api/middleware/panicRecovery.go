package middleware

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"

	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type PanicRecoveryMiddleware struct {
	logger logrus.FieldLogger
}

type PanicRecoveryMiddlewareDependencies struct {
	Logger logrus.FieldLogger
}

func NewPanicRecoveryMiddleware(dep PanicRecoveryMiddlewareDependencies) (*PanicRecoveryMiddleware, error) {
	if dep.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("PanicRecovery", "Logger", "nil")
	}

	logger := dep.Logger.WithField("middleware", "panicRecovery")

	return &PanicRecoveryMiddleware{
		logger: logger,
	}, nil
}

func (mw *PanicRecoveryMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panicErr := NewErrPanicWrapper(err)
				log.Print(panicErr)
				return
			}
		}()

		next.ServeHTTP(w, r)

	})
}
