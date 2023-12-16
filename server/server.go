package server

import (
	"mpmy-product-service/config"
	"mpmy-product-service/server/routes"

	"github.com/gin-gonic/gin"

	"mpmy-product-service/graph"
	"mpmy-product-service/middleware"
	"mpmy-product-service/service"
)

type Server struct {
	loginService           *service.LoginService
	productService         *service.ProductService
	categoryProductService *service.CategoryProductService
	priceScheduleService   *service.PriceScheduleService
	categoryService        *service.CategoryService
	recentSearchService	*service.RecentSearchService
}

func NewServer(loginService *service.LoginService,
	productService *service.ProductService,
	categoryProductService *service.CategoryProductService,
	priceScheduleService *service.PriceScheduleService,
	categoryService *service.CategoryService,
	recentSearchService *service.RecentSearchService) *Server {
	return &Server{
		loginService:           loginService,
		productService:         productService,
		categoryProductService: categoryProductService,
		priceScheduleService:   priceScheduleService,
		categoryService:        categoryService,
		recentSearchService: recentSearchService,
	}
}

func (server *Server) RoutesHandler(appConfig config.AppConfig) *gin.Engine {
	r := gin.Default()
	r.HandleMethodNotAllowed = true
	r.Use(middleware.CORSMiddleware())
	// r.Use(cors.Default())
	r.Use(middleware.ReqBodyMiddleware())
	r.Use(middleware.AuthMiddleware(appConfig))
	r.Use(middleware.GinContextToContextMiddleware())

	// prepare resolvers
	resolvers := &graph.Resolver{
		LoginService:           server.loginService,
		ProductService:         server.productService,
		CategoryProductService: server.categoryProductService,
		PriceScheduleService:   server.priceScheduleService,
		CategoryService:        server.categoryService,
		RecentSearchService:   server.recentSearchService,
	}

	routes.Routes(r, resolvers, appConfig)

	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"errCode": 405, "errMsg": "Method Not Allowed"})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"errCode": 404, "errMsg": "Not Found"})
	})

	return r
}
