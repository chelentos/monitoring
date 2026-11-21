package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseWriterWithLogging struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *ResponseWriterWithLogging) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	if rw.statusCode == 0 {
		rw.statusCode = statusCode
	}
}

func (rw *ResponseWriterWithLogging) Write(data []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	size, err := rw.ResponseWriter.Write(data)
	rw.size += size

	return size, err
}

func WithLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			uri := r.RequestURI
			method := r.Method

			cw := ResponseWriterWithLogging{
				ResponseWriter: w,
				statusCode:     0,
				size:           0,
			}

			next.ServeHTTP(&cw, r)

			duration := time.Since(start)
			if cw.statusCode == 0 {
				cw.statusCode = 200
			}

			logger.Info(
				"server",
				zap.String("URI", uri),
				zap.String("Method", method),
				zap.Int("Status code", cw.statusCode),
				zap.Int("Size", cw.size),
				zap.Duration("Duration", duration),
			)
		})
	}
}
