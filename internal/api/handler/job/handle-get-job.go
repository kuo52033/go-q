package jobHandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) HandleGetJob(c *gin.Context) {
	req := c.MustGet("dto").(*GetJobRequest)

	log.Println(req)

	job, err := h.jobService.GetJobById(c.Request.Context(), req.JobID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, job)
}
