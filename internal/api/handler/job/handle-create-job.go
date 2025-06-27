package jobHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleCreateJob(c *gin.Context) {
	req := c.MustGet("dto").(*CreateJobRequest)

	job, err := h.jobService.CreateJob(c.Request.Context(), req.JobType, req.Payload, req.QueueName, req.MaxAttempts)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, job)
}
