package mnemo

import (
	"fmt"
)

type (
	// Error is a generic error type for the Mnemo package.
	Error[T any] struct {
		Err    error
		Status int
		Logger Logger
		level  LogLevel
	}
)

// NewError returns a new Error instance with a logger
func NewError[T any](msg string, opts ...Opt[Error[T]]) Error[T] {
	e := Error[T]{
		Err:    fmt.Errorf(msg),
		Logger: logger,
	}
	return e
}

// Error implements the error interface.
func (e Error[T]) Error() string {
	return fmt.Sprintf("%v error: %v", new(T), e.Err.Error())
}

// IsStatusError returns true if the error has a status code.
func (e Error[T]) IsStatusError() bool {
	return e.Status != 0
}

// Log logs the error to the logger by log level.
func (e Error[T]) Log() {
	switch e.level {
	case Debug:
		e.Logger.Debug(e.Err.Error())
	case Info:
		e.Logger.Info(e.Err.Error())
	case Warn:
		e.Logger.Warn(e.Err.Error())
	case Fatal:
		e.Logger.Fatal(e.Err.Error())
	case Panic:
		panic(e.Err.Error())
	default:
		e.Logger.Error(e.Err.Error())
	}
}

// WithStatus sets the status code for the error.
func (e Error[T]) WithStatus(status int) Error[T] {
	e.Status = status
	return e
}

// WithLogLevel sets the log level for the error.
func (e Error[T]) WithLogLevel(level LogLevel) Error[T] {
	e.level = level
	return e
}

// IsErrorType reflects Error[T] from an error.
func IsErrorType[T any](err error) (Error[T], bool) {
	t, ok := err.(Error[T])
	return t, ok
}
