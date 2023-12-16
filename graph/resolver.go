package graph

import "mpmy-product-service/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	LoginService           *service.LoginService
	ProductService         *service.ProductService
	CategoryProductService *service.CategoryProductService
	PriceScheduleService   *service.PriceScheduleService
	CategoryService        *service.CategoryService
	RecentSearchService    *service.RecentSearchService
}
