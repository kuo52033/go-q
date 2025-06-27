package jobMiddleware

import (
	"encoding/json"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	jobHandler "github.com/kuo-52033/go-q/internal/api/handler/job"
	"github.com/kuo-52033/go-q/internal/utils/myerror"
	validatorUtil "github.com/kuo-52033/go-q/internal/utils/validator"
)

func ValidateQueuePayload() gin.HandlerFunc {
	return func(c *gin.Context) {

		value, exists := c.Get("dto")
		if !exists {
			c.Abort()
			return
		}

		req := value.(*jobHandler.CreateJobRequest)

		var payloadStruct interface{}

		log.Println("req", req)

		switch req.QueueName {
		case "process_image":
			payloadStruct = &jobHandler.ProcessImagePayload{}
		case "send_email":
			payloadStruct = &jobHandler.SendEmailPayload{}
		case "generate_report":
			payloadStruct = &jobHandler.GenerateReportPayload{}
		default:
			return
		}

		
		payloadBytes, err := json.Marshal(req.Payload)
		if err != nil {
			c.Error(myerror.RequestValidationError(myerror.REQUEST_VALIDATION_ERROR, myerror.WithError(err)))
			c.Abort()
			return
		}
		
		err = json.Unmarshal(payloadBytes, payloadStruct)
		if err != nil {
			c.Error(myerror.RequestValidationError(myerror.REQUEST_VALIDATION_ERROR, myerror.WithError(err)))
			c.Abort()
			return
		}
		
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			if err := v.Struct(payloadStruct); err != nil {
				validateError := validatorUtil.FormateValidationError(err)
				appErr := myerror.RequestValidationError(
					myerror.REQUEST_VALIDATION_ERROR,
					func(e *myerror.MyError) {
						if validateError != nil {
							e.Extra = validateError
						} else {
							e.Extra = map[string]interface{}{
								"error": err.Error(),
							}
						}
					},
				)
				c.Error(appErr)				
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
