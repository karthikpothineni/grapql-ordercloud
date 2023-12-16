package utils

import (
	"github.com/gin-gonic/gin"
	"mpmy-product-service/constants"
)

func serverError(c *gin.Context, err string) {
	c.JSON(401, gin.H{
		"Success": false,
		"message": "Internal server error: " + err,
	})
	c.Abort()
}
func HandlePanic(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			Logger(constants.DefaultLogLevel, err)
			serverError(c, err.(error).Error())
			return
		}
	}()
}
