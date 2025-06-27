package common

import (
	"reflect"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/utils/myerror"
	"github.com/kuo-52033/go-q/internal/utils/validator"
)

// hasTag checks if the type has the given tag
func hasTag(t reflect.Type, tag string) bool {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if _, ok := field.Tag.Lookup(tag); ok {
			return true
		}

		if field.Type.Kind() == reflect.Struct {
			if hasTag(field.Type, tag) {
				return true
			}
		}
	}
	return false
}

func Validate(dtoType interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(dtoType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("validate middleware only support struct")
	}

	hasUri := hasTag(t, "uri")
	hasQuery := hasTag(t, "form")
	hasJson := hasTag(t, "json")

	return func(c *gin.Context) {
		dto := reflect.New(t).Interface()

		var finalErr error

		if hasUri {
			if err := c.ShouldBindUri(dto); err != nil {
				finalErr = err
			}
		}

		if finalErr == nil && hasQuery {
			if err := c.ShouldBindQuery(dto); err != nil {
				finalErr = err
			}
		}

		if finalErr == nil && hasJson && c.Request.ContentLength > 0 {
			if err := c.ShouldBindJSON(dto); err != nil {
				finalErr = err
			}
		}

		if finalErr != nil {
			log.Printf("validate error: %v", finalErr)
			validateError := validator.FormateValidationError(finalErr)
			appErr := myerror.RequestValidationError(
				myerror.REQUEST_VALIDATION_ERROR,
				func(e *myerror.MyError) {
					if validateError != nil {
						e.Extra = validateError
					} else {
						e.Extra = map[string]interface{}{
							"error": finalErr.Error(),
						}
					}
				},
			)
			c.Error(appErr)
			c.Abort()
			return
		}

		c.Set("dto", dto)
		c.Next()
	}
}
