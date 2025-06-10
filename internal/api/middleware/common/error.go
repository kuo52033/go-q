package common

import (
	"errors"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/kuo-52033/go-q/internal/error"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var myError *error.Error
			if errors.As(err, &myError) {
				c.JSON(myError.HTTPStatus, gin.H{
					"code": myError.Code,
					"message": myError.Error(),
					"extra": myError.Extra,
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": "000",
					"message": "Internal Server Error",
					"extra": map[string]interface{}{
						"error": err.Error(),
					},
				})
			}
			c.Abort()
		}
	}
}
