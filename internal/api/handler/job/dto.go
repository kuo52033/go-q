package jobHandler

import "github.com/kuo-52033/go-q/internal/model"

type ProcessImagePayload struct {
	ImageURL string `json:"image_url" binding:"required"`
}

type SendEmailPayload struct {
	To string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body string `json:"body" binding:"required"`
}
	
type GenerateReportPayload struct {
	ReportType string `json:"report_type" binding:"required,oneof=daily weekly monthly"`
}

type CreateJobRequest struct {
	JobType string `json:"job_type" binding:"required"`
	Payload model.JobPayload `json:"payload" binding:"required"`
	QueueName string `json:"queue_name" binding:"required,oneof=process_image send_email generate_report"`
	MaxAttempts int `json:"max_attempts" binding:"required,min=1"`
}

type GetJobRequest struct {
	JobID string `uri:"job_id" binding:"required"`
}
