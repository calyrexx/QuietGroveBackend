package middleware

import (
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/errorspkg"
	"github.com/calyrexx/QuietGrooveBackend/internal/pkg/utils"
	"github.com/calyrexx/zeroslog"
	"github.com/gorilla/mux"
	"log/slog"
	"time"

	"net/http"
)

type PanicRecoveryMiddleware struct {
	logger *slog.Logger
}

type PanicRecoveryMiddlewareDependencies struct {
	Logger *slog.Logger
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
				mw.logger.Error("got panic", zeroslog.ErrorKey, panicErr)
				utils.WriteError(w, http.StatusInternalServerError, errorspkg.ErrInternalService)
				return
			}
		}()
		startTime := time.Now()

		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		mw.logger.Info("HTTP request started", "path", path, "method", r.Method)

		srw := &statusCapturingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(srw, r)

		duration := time.Since(startTime)

		statusCode := srw.Status()

		switch statusCode {
		case http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusForbidden,
			http.StatusInternalServerError:
			mw.logger.Error("HTTP request failed", "path", path, "method", r.Method, "duration", duration, "status", statusCode)
		default:
			mw.logger.Info("HTTP request completed", "path", path, "method", r.Method, "duration", duration, "status", statusCode)
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
