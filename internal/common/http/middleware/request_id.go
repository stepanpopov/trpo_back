package middleware

import (
	"net/http"

	"github.com/google/uuid"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
)

// SetReqId is a middleware that generates a unique request ID for each HTTP request.
// The request ID is added to the request context, allowing it to be used throughout the request lifecycle.
//
// Behavior:
//   - Generates a new UUID for each request.
//   - Wraps the request with the generated request ID.
//
// Returns:
//   - A middleware function that wraps an http.Handler and adds a unique request ID to the request context.
func SetReqId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := uuid.New()
		next.ServeHTTP(w, commonHttp.WrapReqID(r, reqId.ID()))
	})
}
