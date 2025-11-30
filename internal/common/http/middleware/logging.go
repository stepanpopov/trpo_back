// Package middleware provides HTTP middleware for common tasks such as logging and authentication.
package middleware

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

const realIPHeaderName = "X-Real-IP"

// Logging is a middleware that logs HTTP requests and their response times.
// It uses the provided logger to log the HTTP method, URL path, client IP address, and the time taken to process the request.
// If the "X-Real-IP" header is present, it logs the real client IP address; otherwise, it falls back to the remote address.
//
// Parameters:
//   - logger: An instance of logger.Logger used to log the request details.
//
// Returns:
//   - A middleware function that wraps an http.Handler and logs request details.
func Logging(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				respTime := time.Since(start)
				realIP := r.Header.Get(realIPHeaderName)
				if realIP != "" {
					logger.InfofReqID(r.Context(), "%s %s from RealIP: %s IP: %s - %s",
						r.Method, r.URL.Path, realIP, r.RemoteAddr, respTime.String())
				} else {
					logger.InfofReqID(r.Context(), "%s %s from IP: %s - %s",
						r.Method, r.URL.Path, r.RemoteAddr, respTime.String())
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
