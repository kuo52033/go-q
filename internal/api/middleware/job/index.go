package jobMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/middleware/common"
	jobHandler "github.com/kuo-52033/go-q/internal/api/handler/job"
	jobMiddleware "github.com/kuo-52033/go-q/internal/api/middleware/job/before-create-job-request"
)

type Middleware struct {
	CreateJobMiddleware gin.HandlersChain
	GetJobMiddleware gin.HandlersChain
}

func NewMiddleware() *Middleware {
	return &Middleware{
		CreateJobMiddleware: gin.HandlersChain{
			common.Validate(jobHandler.CreateJobRequest{}),
			jobMiddleware.ValidateQueuePayload(),
		},
		GetJobMiddleware: gin.HandlersChain{
			common.Validate(jobHandler.GetJobRequest{}),
		},
	}
}
