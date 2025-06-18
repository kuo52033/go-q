package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/handler/job"
	"github.com/kuo-52033/go-q/internal/service"
	"github.com/kuo-52033/go-q/internal/store"
	"github.com/redis/go-redis/v9"
)

func SetupJobModule(router *gin.RouterGroup, rdb *redis.Client) {
	jobStore := store.NewRedisJobStore(rdb)
	jobService := service.NewJobService(jobStore)
	jobHandler := job.NewHandler(jobService)

	jobGroup := router.Group("/jobs")
	{
		jobGroup.POST(
			"/", 
			jobHandler.HandleCreateJob,
		)
		jobGroup.GET(
			"/:job_id",
			jobHandler.HandleGetJob,
		)
	}
}
