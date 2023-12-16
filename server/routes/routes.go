package routes

import (
	"mpmy-product-service/config"
	"mpmy-product-service/graph"

	"github.com/gin-gonic/gin"
)

// Routes ...
func Routes(r *gin.Engine, resolvers *graph.Resolver, appConfig config.AppConfig) {
	Query(r, resolvers, appConfig)
	QueryRoot(r, appConfig)
	Health(r)
}
