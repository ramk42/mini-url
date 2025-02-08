package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ramk42/mini-url/internal/infra/logging"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)
		requestID := middleware.GetReqID(r.Context())
		duration := time.Since(start)

		logEvent := logging.Logger.Info().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", recorder.statusCode).
			Dur("duration", duration).
			Str("ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent())

		if logging.Logger.GetLevel() == zerolog.DebugLevel {
			logEvent = logEvent.
				Interface("request_headers", r.Header).
				Interface("response_headers", recorder.Header())
		}

		switch {
		case recorder.statusCode >= 500:
			logEvent.Send()
		case recorder.statusCode >= 400:
			logEvent.Send()
		default:
			logEvent.Send()
		}
	})
}
