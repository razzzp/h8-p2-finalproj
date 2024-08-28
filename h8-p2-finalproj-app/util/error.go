package util

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("status code: %d;message: %s;detail: %s", ae.StatusCode, ae.Message, ae.Detail)
}

func NewAppError(status int, message string, detail string) *AppError {
	return &AppError{
		StatusCode: status,
		Message:    message,
		Detail:     detail,
	}
}
func ErrorHandler(e error, c echo.Context) {
	if httperror, ok := e.(*echo.HTTPError); ok {
		c.JSON(httperror.Code, map[string]any{
			"message": httperror.Message,
		})
	} else if apperror, ok := e.(*AppError); ok {
		c.JSON(apperror.StatusCode, apperror)
	} else {
		c.JSON(http.StatusInternalServerError, e.Error())
	}
}
