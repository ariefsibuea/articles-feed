package handler

import (
	"fmt"

	_errors "github.com/ariefsibuea/articles-feed/internal/pkg/errors"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"-"`
}

type Meta struct {
	Page       int32 `json:"page,omitempty"`
	PageSize   int32 `json:"pageSize,omitempty"`
	TotalItems int32 `json:"totalItems,omitempty"`
}

func Success(c echo.Context, statusCode int, data interface{}, meta *Meta) error {
	if meta == nil {
		meta = &Meta{}
	}

	return c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code    int
			message string
			details string
		)

		code = _errors.GetErrorCode(err)

		message = err.Error()
		if c.Echo().Debug {
			details = fmt.Sprintf("%+v", err)
		}

		if e, ok := err.(*echo.HTTPError); ok {
			code = e.Code
			message = fmt.Sprintf("%+v", e.Message)
		}

		c.JSON(code, Response{
			Success: false,
			Error: &Error{
				Code:    code,
				Message: message,
				Details: details,
			},
		})
	}
}
