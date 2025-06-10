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
)
