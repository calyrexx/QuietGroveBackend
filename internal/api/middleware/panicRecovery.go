package middleware

import (
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/Calyr3x/QuietGrooveBackend/internal/pkg/utils"
	"github.com/gorilla/mux"
	"time"

	"net/http"

	"github.com/sirupsen/logrus"
)

type PanicRecoveryMiddleware struct {
	logger logrus.FieldLogger
}

type PanicRecoveryMiddlewareDependencies struct {
	Logger logrus.FieldLogger
}

func NewPanicRecoveryMiddleware(d PanicRecoveryMiddlewareDependencies) (*PanicRecoveryMiddleware, error) {
	if d.Logger == nil {
		return nil, errorspkg.NewErrConstructorDependencies("PanicRecovery", "Logger", "nil")
	}

	return &PanicRecoveryMiddleware{
		logger: d.Logger,
	}, nil
}

func (mw *PanicRecoveryMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panicErr := errorspkg.NewErrPanicWrapper(err)
				mw.logger.Error(panicErr)
				utils.WriteError(w, http.StatusInternalServerError, errorspkg.ErrInternalService)
				return
			}
		}()
		startTime := time.Now()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		mw.logger.WithFields(logrus.Fields{
			"path":   path,
			"method": r.Method,
		}).Info("HTTP request started")

		srw := &statusCapturingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(srw, r)

		duration := time.Since(startTime)

		statusCode := srw.Status()

		entry := mw.logger.WithFields(logrus.Fields{
			"path":     path,
			"method":   r.Method,
			"duration": duration,
			"status":   statusCode,
		})
		switch statusCode {
		case http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusForbidden,
			http.StatusInternalServerError:
			entry.Error("HTTP request failed")
		default:
			entry.Info("HTTP request completed")
		}
	})
}

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (w *statusCapturingResponseWriter) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.status = statusCode
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCapturingResponseWriter) Status() int {
	if w.wroteHeader {
		return w.status
	}
	// если статус не был явно установлен, по умолчанию 200
	return http.StatusOK
}
