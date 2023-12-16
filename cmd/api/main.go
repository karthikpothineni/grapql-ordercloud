package main

import (
	"context"
	"github.com/dgraph-io/ristretto"
	"mpmy-product-service/client/cache"
	"mpmy-product-service/client/db"

	"mpmy-product-service/config"
	"mpmy-product-service/constants"
	"mpmy-product-service/server"
	"mpmy-product-service/service"
	"mpmy-product-service/utils"
)

func main() {
	// init config
	appConfig := config.Init()
	name := appConfig.General.Name
	port := appConfig.General.Port

	// init db
	dbClient, err := db.Init(appConfig.DB)
	if err != nil {
		panic(err)
	}

	// init cache
	rCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,
		MaxCost:     int64(500 * 1024 * 1024),
		BufferItems: 64,
		Metrics:     false,
	})
	if err != nil {
		panic(err)
	}
	cacheClient := cache.New(rCache)

	// init services
	loginService := service.NewLoginService()
	productService := service.NewProductService(dbClient, appConfig.DB, appConfig.OrderCloud, cacheClient)
	categoryProductService := service.NewCategoryProductService(appConfig.OrderCloud)
	priceScheduleService := service.NewPriceScheduleService()
	categoryService := service.NewCategoryService(appConfig.OrderCloud)
	recentSearchService := service.NewRecentSearchService(dbClient, appConfig.DB)

	// fetch order cloud access token asynchronously
	go func() {
		loginService.StartAccessTokenFetcher(context.Background(), appConfig.OrderCloud)
	}()

	// start trending products processor
	go func() {
		productService.StartTrendingProductsProcessor(context.Background())
	}()

	// init server
	appServer := server.NewServer(loginService, productService, categoryProductService, priceScheduleService, categoryService, recentSearchService)
	r := appServer.RoutesHandler(appConfig)

	// start server
	utils.Logger(constants.DefaultLogLevel, name+" start listening at http://localhost:"+port)
	utils.Logger(constants.DefaultLogLevel, "==> 🚀 %s listening at %s\n"+name+port)
	err = r.Run(":" + port)
	if err != nil {
		return
	}

}
