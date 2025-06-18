package job

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) HandleGetJob(c *gin.Context) {
	jobID := c.Param("job_id")

	job, err := h.jobService.GetJobById(c.Request.Context(), jobID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, job)
}
