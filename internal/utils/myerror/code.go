package myerror

type ErrorCode struct{
	Code string
	Message string
}

var (
	JOB_CREATE_FAILED = ErrorCode{
		Code: "001",
		Message: "Failed to create job",
	}
	JOB_GET_FAILED = ErrorCode{
		Code: "002",
		Message: "Failed to get job",
	}
	REQUEST_VALIDATION_ERROR = ErrorCode{
		Code: "003",
		Message: "Request validation error",
	}
)
