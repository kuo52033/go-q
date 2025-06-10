package myerror

import "net/http"

type MyError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Extra      map[string]interface{} `json:"extra"`
}

func (e *MyError) Error() string {
	return e.Message
}

func NotFound(e ErrorCode, extra map[string]interface{}) *MyError {
	return &MyError{
		HTTPStatus: http.StatusNotFound,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func RequestValidationError(e ErrorCode, extra map[string]interface{}) *MyError {
	return &MyError{
		HTTPStatus: http.StatusBadRequest,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func Forbidden(e ErrorCode, extra map[string]interface{}) *MyError {
	return &MyError{
		HTTPStatus: http.StatusForbidden,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func InternalServerError(e ErrorCode, extra map[string]interface{}) *MyError {
	return &MyError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func Unauthorized(e ErrorCode, extra map[string]interface{}) *MyError {
	return &MyError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}
