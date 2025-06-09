package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/model"
)

type CreateJobRequest struct {
	JobType string `json:"job_type"`
	Payload model.JobPayload `json:"payload"`
	QueueName string `json:"queue_name"`
	MaxAttempts int `json:"max_attempts"`
}

func (h *Handler) HandleCreateJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobService.CreateJob(c.Request.Context(), req.JobType, req.Payload, req.QueueName, req.MaxAttempts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}
