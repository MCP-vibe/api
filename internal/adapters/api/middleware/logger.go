package middleware

import (
	"api/internal/adapters/api/logging"
	"api/internal/adapters/logger"
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"
)

type Logger struct {
	log logger.Logger
}

func NewLogger(log logger.Logger) Logger {
	return Logger{log: log}
}

func (l Logger) Execute(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	const (
		logKey      = "logger_middleware"
		requestKey  = "api_request"
		responseKey = "api_response"
	)

	body, err := getRequestPayload(r)
	if err != nil {
		logging.NewError(l.log, err, logKey, http.StatusBadRequest).Log("error when getting payload")
		return
	}

	l.log.WithFields(logger.Fields{
		"key":         requestKey,
		"payload":     body,
		"url":         r.URL.Path,
		"http_method": r.Method,
	}).Infof("started handling request")

	next.ServeHTTP(w, r)

	end := time.Since(start).Seconds()
	res := w.(negroni.ResponseWriter)
	l.log.WithFields(logger.Fields{
		"key":           responseKey,
		"url":           r.URL.Path,
		"http_method":   r.Method,
		"http_status":   res.Status(),
		"response_time": end,
	}).Infof("completed handling request")
}
func getRequestPayload(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", errors.New("body not defined")
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("error read body")
	}

	r.Body = io.NopCloser(bytes.NewBuffer(payload))

	return strings.TrimSpace(string(payload)), nil
}
