package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/api/handler/job"
)

type RouteConfig struct {
	JobHandler *job.Handler
}

func SetupRoutes(router *gin.Engine, config *RouteConfig) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")

	SetupJobRoutes(api, config.JobHandler)

}
