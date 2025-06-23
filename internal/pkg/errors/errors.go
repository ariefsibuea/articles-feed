package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidSearchPath = errors.New("invalid search path")

	ErrAuthorNotFound = NotFoundErrorf("author not found")
)

type CustomError interface {
	error
	Code() int
}

func GetErrorCode(err error) int {
	var errCustom CustomError
	if errors.As(err, &errCustom) {
		return errCustom.Code()
	}

	return http.StatusInternalServerError
}

type BadRequestError struct {
	statusCode int
	message    string
}

func (e *BadRequestError) Code() int {
	return e.statusCode
}

func (e *BadRequestError) Error() string {
	return e.message
}

func BadRequestErrorf(format string, args ...interface{}) CustomError {
	return &BadRequestError{
		statusCode: http.StatusBadRequest,
		message:    fmt.Sprintf(format, args...),
	}
}

type UnauthorizedError struct {
	statusCode int
	message    string
}

func (e *UnauthorizedError) Code() int {
	return e.statusCode
}

func (e *UnauthorizedError) Error() string {
	return e.message
}

func UnauthorizedErrorf(format string, args ...interface{}) CustomError {
	return &UnauthorizedError{
		statusCode: http.StatusUnauthorized,
		message:    fmt.Sprintf(format, args...),
	}
}

type NotFoundError struct {
	statusCode int
	message    string
}

func (e *NotFoundError) Code() int {
	return e.statusCode
}

func (e *NotFoundError) Error() string {
	return e.message
}

func NotFoundErrorf(format string, args ...interface{}) CustomError {
	return &NotFoundError{
		statusCode: http.StatusNotFound,
		message:    fmt.Sprintf(format, args...),
	}
}
