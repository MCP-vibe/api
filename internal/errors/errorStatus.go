package errors

import "errors"

type ErrorStatus struct {
	StatusCode int
	Err        error
}

func NewError(statusCode int, err string) *ErrorStatus {
	return &ErrorStatus{
		StatusCode: statusCode,
		Err:        errors.New(err),
	}
}

func (r *ErrorStatus) Error() string {
	return r.Err.Error()
}
