package middleware

import (
	"net/http"
	"runtime/debug"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger"
)

// Panic is a middleware that recovers from panics in the HTTP handler chain and logs the error.
// It ensures that the server does not crash due to unexpected panics and responds with a generic error message.
//
// Parameters:
//   - logger: An instance of logger.Logger used to log the panic details and stack trace.
//
// Behavior:
//   - If a panic occurs, it logs the error message and stack trace using the provided logger.
//   - Sends a generic "server unknown error" response with HTTP status 500.
//
// Returns:
//   - A middleware function that wraps an http.Handler and recovers from panics.
func Panic(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorf("PANIC (recovered): %s\n stacktrace:\n%s", err, string(debug.Stack()))
					commonHttp.ErrorResponse(w, r, "server unknown error", http.StatusInternalServerError, logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
