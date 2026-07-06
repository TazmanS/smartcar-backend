package middleware

import (
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logger.Info(
			"In",
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}
