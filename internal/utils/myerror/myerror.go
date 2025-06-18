package myerror

import "net/http"

type MyError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Extra      map[string]interface{} `json:"extra"`
	Err        error                  `json:"-"`
}

type Option func(*MyError)

func WithError(err error) Option {
	return func(e *MyError) {
		if e.Extra == nil {
			e.Extra = make(map[string]interface{})
		}
		e.Extra["error"] = err.Error()
		e.Err = err
	}
}

func (e *MyError) Error() string {
	return e.Message
}

func (e *MyError) Unwrap() error {
	return e.Err
}

func NotFound(e ErrorCode, opts ...Option) *MyError {
	err := &MyError{
		HTTPStatus: http.StatusNotFound,
		Code:       e.Code,
		Message:    e.Message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func RequestValidationError(e ErrorCode, opts ...Option) *MyError {
	err := &MyError{
		HTTPStatus: http.StatusBadRequest,
		Code:       e.Code,
		Message:    e.Message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func Forbidden(e ErrorCode, opts ...Option) *MyError {
	err := &MyError{
		HTTPStatus: http.StatusForbidden,
		Code:       e.Code,
		Message:    e.Message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func InternalServerError(e ErrorCode, opts ...Option) *MyError {
	err := &MyError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       e.Code,
		Message:    e.Message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}

func Unauthorized(e ErrorCode, opts ...Option) *MyError {
	err := &MyError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       e.Code,
		Message:    e.Message,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
