// Package logger provides an interface and implementation for logging messages with optional request ID support.
// It allows logging of error and informational messages, with or without formatting, and supports associating logs with request IDs.
package logger

import (
	"context"
)

//go:generate mockgen -source=logger.go -destination=mocks/mock.go

// ReqIDGetter is a function type that retrieves a request ID from the given context.
// It returns the request ID as a uint32 and an error if the request ID cannot be retrieved.
type ReqIDGetter func(ctx context.Context) (uint32, error)

// Logger defines an interface for logging messages with support for request IDs.
// It provides methods for logging error and informational messages, with or without formatting.
//
// Methods:
//   - Error: Logs an error message.
//   - Errorf: Logs a formatted error message.
//   - Info: Logs an informational message.
//   - Infof: Logs a formatted informational message.
//   - ErrorReqID: Logs an error message with a request ID from the context.
//   - ErrorfReqID: Logs a formatted error message with a request ID from the context.
//   - InfoReqID: Logs an informational message with a request ID from the context.
//   - InfofReqID: Logs a formatted informational message with a request ID from the context.
type Logger interface {
	Error(msg string)
	Errorf(format string, a ...any)
	Info(msg string)
	Infof(format string, a ...any)

	ErrorReqID(ctx context.Context, msg string)
	ErrorfReqID(ctx context.Context, format string, a ...any)
	InfoReqID(ctx context.Context, msg string)
	InfofReqID(ctx context.Context, format string, a ...any)
}

// NewLogger creates a new instance of Logger using the provided ReqIDGetter function.
// The ReqIDGetter is used to retrieve request IDs from the context.
//
// Parameters:
//   - reqIdGetter: A function that retrieves request IDs from the context.
//
// Returns:
//   - A Logger instance.
//   - An error if the logger initialization fails.
func NewLogger(reqIdGetter ReqIDGetter) (Logger, error) {
	return NewFLogger(reqIdGetter)
}
