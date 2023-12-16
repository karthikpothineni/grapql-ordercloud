package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mpmy-product-service/constants"
	"mpmy-product-service/utils"
)

func ReqBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.HandlePanic(c)
		ByteBody, err := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(ByteBody))
		if ByteBody != nil {
			jsonData := string(ByteBody)
			utils.Logger(constants.DefaultLogLevel, "/*jsonDatajsonData", jsonData, "jsonDatajsonData*/")
			if err != nil {
				utils.Logger(constants.DefaultLogLevel, "ReqBodyMiddleware err", err)
			}
			if jsonData != "" {
				c.Set(constants.RequestBody, jsonData)
			}
		}
		c.Next()
	}
}
