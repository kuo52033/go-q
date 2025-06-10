package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/handler/job"
	"github.com/kuo-52033/go-q/internal/db"
	"github.com/kuo-52033/go-q/internal/service"
)

func SetupJobModule(router *gin.RouterGroup, rdb db.RedisClient) {
	jobService := service.NewJobService(rdb)
	jobHandler := job.NewHandler(jobService)

	jobGroup := router.Group("/jobs")
	{
		jobGroup.POST(
			"/", 
			jobHandler.HandleCreateJob,
		)
	}
}
