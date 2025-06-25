package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/model"
)

type CreateJobRequest struct {
	JobType string `json:"JobType"`
	Payload model.JobPayload `json:"Payload"`
	QueueName string `json:"QueueName"`
	MaxAttempts int `json:"MaxAttempts"`
}

func (h *Handler) HandleCreateJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	job, err := h.jobService.CreateJob(c.Request.Context(), req.JobType, req.Payload, req.QueueName, req.MaxAttempts)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, job)
}
