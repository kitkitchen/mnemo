package mnemo

import (
	"fmt"
)

type (
	Error[T any] struct {
		Err    error
		Status int
		Logger Logger
		level  LogLevel
	}
)

func NewError[T any](msg string, opts ...Opt[Error[T]]) Error[T] {
	e := Error[T]{
		Err:    fmt.Errorf(msg),
		Logger: logger,
	}
	return e
}

func (e Error[T]) Error() string {
	return fmt.Sprintf("%v error: %v", new(T), e.Err.Error())
}

func (e Error[T]) IsStatusError() bool {
	return e.Status != 0
}

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

func (e Error[T]) WithStatus(status int) Error[T] {
	e.Status = status
	return e
}

func (e Error[T]) WithLogLevel(level LogLevel) Error[T] {
	e.level = level
	return e
}

func IsErrorType[T any](err error) bool {
	_, ok := err.(Error[T])
	return ok
}
