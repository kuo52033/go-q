package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/handler/job"
)

func SetupJobRoutes(router *gin.RouterGroup, jobHandler *job.Handler) {
	jobGroup := router.Group("/jobs")
	{
		jobGroup.POST(
			"/", 
			jobHandler.HandleCreateJob,
		)
	}
}
