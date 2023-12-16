package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"mpmy-product-service/config"
	"mpmy-product-service/graph"
	"mpmy-product-service/graph/generated"
)

func PlaygroundHandler(appConfig config.AppConfig) gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
func PlaygroundHandlerWithToken(appConfig config.AppConfig) gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/with_token/"+appConfig.General.StaticToken+"/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GraphqlHandler(resolvers *graph.Resolver) gin.HandlerFunc {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
