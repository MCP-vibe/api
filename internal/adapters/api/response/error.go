package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"api/internal/adapters/api/logging"
	"api/internal/adapters/logger"
	errorStatus "api/internal/errors"
)

var ErrInvalidInput = errors.New("invalid_input")

type Error struct {
	statusCode int
	Errors     []string `json:"errors"`
}

func NewError(err error, status int) *Error {
	return &Error{
		statusCode: status,
		Errors:     []string{err.Error()},
	}
}

func NewErrorMessage(messages []string, status int) *Error {
	return &Error{
		statusCode: status,
		Errors:     messages,
	}
}

func (e Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	return json.NewEncoder(w).Encode(e)
}

func NewErrorTokenNotFound(log logger.Logger, logKey string, w http.ResponseWriter) {
	err := errors.New("could not get token from context")
	logging.NewError(log, err, logKey, http.StatusUnauthorized).Log("token not found in context")

	NewError(err, http.StatusUnauthorized).Send(w)
}

func NewErrorWithErrorStatus(err error, w http.ResponseWriter, log logger.Logger, logKey, message string) {
	re, ok := err.(*errorStatus.ErrorStatus)
	var statusCode int
	if ok {
		statusCode = re.StatusCode
	} else {
		statusCode = http.StatusInternalServerError
	}
	NewError(err, statusCode).Send(w)

	logMsg := message

	logging.NewError(log, err, logKey, statusCode).Log(logMsg)
}
