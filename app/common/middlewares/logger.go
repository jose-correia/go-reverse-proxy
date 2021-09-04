package middlewares

import (
	"net/http"

	"github.com/go-kit/kit/log"
)

func preRecord(r *http.Request) []interface{} {
	return []interface{}{
		"transport", "http",
		"direction", "incoming",
		"middleware", "before",
		"uri", r.RequestURI,
		"method", r.Method,
	}
}

func postRecord(r *http.Request) []interface{} {
	return []interface{}{
		"transport", "http",
		"direction", "incoming",
		"middleware", "after",
		"uri", r.RequestURI,
		"method", r.Method,
	}
}

func Logger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer logger.Log(postRecord(r)...)
			logger.Log(preRecord(r)...)
			next.ServeHTTP(w, r)
		})
	}
}
