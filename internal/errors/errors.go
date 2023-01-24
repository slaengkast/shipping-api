package errors

import (
	errs "errors"
)

type ErrorType int

const (
	ErrorNotFound ErrorType = iota
	ErrorConflict
	ErrorUnknown
	ErrorInternal
	ErrorInput
)

type APIError struct {
	t   ErrorType
	err error
}

func FromError(err error, t ErrorType) error {
	return APIError{t: t, err: err}
}

func FromMessage(msg string, t ErrorType) error {
	return APIError{t: t, err: errs.New(msg)}
}

func (e APIError) Error() string {
	return e.err.Error()
}

func (e APIError) GetType() ErrorType {
	return e.t
}
