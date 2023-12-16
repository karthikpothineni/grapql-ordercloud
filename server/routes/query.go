package routes

import (
	"mpmy-product-service/config"
	"mpmy-product-service/graph"
	"mpmy-product-service/server/handler"

	"github.com/gin-gonic/gin"
)

// QueryRoot ...
func QueryRoot(r *gin.Engine, appConfig config.AppConfig) {
	r.GET("/playground", handler.PlaygroundHandler(appConfig))
	if appConfig.General.ByPassSecurity {
		r.GET("/with_token/:jwt_token", handler.PlaygroundHandlerWithToken(appConfig))
	}
}

// Query ...
func Query(r *gin.Engine, resolvers *graph.Resolver, appConfig config.AppConfig) {
	r.POST("/query", handler.GraphqlHandler(resolvers))
	if appConfig.General.ByPassSecurity {
		r.GET("/with_token/:jwt_token/query", handler.GraphqlHandler(resolvers))
		r.POST("/with_token/:jwt_token/query", handler.GraphqlHandler(resolvers))
	}
}
