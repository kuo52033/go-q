package job

import "github.com/kuo-52033/go-q/internal/service"

type Handler struct {
	jobService service.JobService
}

func NewHandler(jobService service.JobService) *Handler {
	return &Handler{jobService: jobService}
}	
