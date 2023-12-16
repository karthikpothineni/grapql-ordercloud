// Package handler - maintains handlers for routes
package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Health check handler
func Health(c *gin.Context) {
	html := fmt.Sprintf(`ok`)
	c.Data(200, "text/html; charset=utf-8", []byte(html))
}
