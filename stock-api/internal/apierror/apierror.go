package apierror

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (err APIError) Error() string {
	return fmt.Sprintf("status: %d, code: %s, message: %s", err.Status, err.Code, err.Message)
}

func NewStatusBadRequestError(message string) APIError {
	return NewAPIError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), message)
}

func NewStatusUnauthorizedError(message string) APIError {
	return NewAPIError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), message)
}

func NewStatusInternalServerError(message string) APIError {
	return NewAPIError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), message)
}

func NewStatusNotFoundError(message string) APIError {
	return NewAPIError(http.StatusNotFound, http.StatusText(http.StatusNotFound), message)
}

func NewAPIError(status int, code, message string) APIError {
	return APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}
