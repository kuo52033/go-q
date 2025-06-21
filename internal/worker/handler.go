package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kuo-52033/go-q/internal/model"
)

// EmailHandler 處理郵件發送任務
type EmailHandler struct{}

func (h *EmailHandler) GetJobType() string {
	return "send_email"
}

func (h *EmailHandler) HandleJob(ctx context.Context, job *model.Job) error {
	log.Printf("Sending email for job %s", job.ID)
	
	// 從 payload 獲取郵件參數
	to, ok := job.Payload["to"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'to' field")
	}
	
	subject, ok := job.Payload["subject"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'subject' field")
	}
	
	// 模擬郵件發送（這裡可以接入實際的郵件服務）
	time.Sleep(2 * time.Second)
	
	log.Printf("Email sent to %s with subject: %s", to, subject)
	return nil
}

// ImageHandler 處理圖片處理任務
type ImageHandler struct{}

func (h *ImageHandler) GetJobType() string {
	return "process_image"
}

func (h *ImageHandler) HandleJob(ctx context.Context, job *model.Job) error {
	log.Printf("Processing image for job %s", job.ID)
	
	imageUrl, ok := job.Payload["image_url"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'image_url' field")
	}
	
	operation, ok := job.Payload["operation"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'operation' field")
	}

	// 模擬圖片處理
	time.Sleep(5 * time.Second)
	
	log.Printf("Image processed: %s with operation: %s", imageUrl, operation)
	return nil
}

// ReportHandler 處理報表生成任務
type ReportHandler struct{}

func (h *ReportHandler) GetJobType() string {
	return "generate_report"
}

func (h *ReportHandler) HandleJob(ctx context.Context, job *model.Job) error {
	log.Printf("Generating report for job %s", job.ID)
	
	reportType, ok := job.Payload["report_type"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'report_type' field")
	}
	
	dateRange, ok := job.Payload["date_range"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid 'date_range' field")
	}

	// 模擬報表生成
	time.Sleep(10 * time.Second)
	
	log.Printf("Report generated: %s for date range: %s", reportType, dateRange)
	return nil
}
