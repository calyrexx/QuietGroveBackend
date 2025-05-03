package middleware

import (
	"errors"
	"fmt"
)

var (
	ErrPanicResponse = errors.New("service got error, try again later")
)

type ErrPanicWrapper struct {
	err interface{}
}

func (e ErrPanicWrapper) Error() string {
	return fmt.Sprintf("got panic: %v", e.err)
}

func NewErrPanicWrapper(err interface{}) error {
	return &ErrPanicWrapper{err: err}
}
