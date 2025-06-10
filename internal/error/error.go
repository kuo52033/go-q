package error

import "net/http"

type Error struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Extra      map[string]interface{} `json:"extra"`
}

func (e *Error) Error() string {
	return e.Message
}

func NotFound(e ErrorCode, extra map[string]interface{}) *Error {
	return &Error{
		HTTPStatus: http.StatusNotFound,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func RequestValidationError(e ErrorCode, extra map[string]interface{}) *Error {
	return &Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func Forbidden(e ErrorCode, extra map[string]interface{}) *Error {
	return &Error{
		HTTPStatus: http.StatusForbidden,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func InternalServerError(e ErrorCode, extra map[string]interface{}) *Error {
	return &Error{
		HTTPStatus: http.StatusInternalServerError,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}

func Unauthorized(e ErrorCode, extra map[string]interface{}) *Error {
	return &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       e.Code,
		Message:    e.Message,
		Extra:      extra,
	}
}
