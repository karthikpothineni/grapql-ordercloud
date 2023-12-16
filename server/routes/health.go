// Package routes ...
package routes

import (
	"mpmy-product-service/server/handler"

	"github.com/gin-gonic/gin"
)

// Health  ...
func Health(r *gin.Engine) {
	r.GET("/health", handler.Health)
}
