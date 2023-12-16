package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"mpmy-product-service/client/cache"
	"mpmy-product-service/config"
	"mpmy-product-service/graph/model"
	"mpmy-product-service/repository"

	"github.com/jmoiron/sqlx"
)

const (
	// TrendingProductCacheKey is the key for trending products cache
	TrendingProductCacheKey = "trending_products"
)

type ProductService struct {
	cacheClient          *cache.Cache
	priceScheduleService IPriceScheduleService
	productRepo          repository.IProductRepository
	categoryProductRepo  repository.ICategoryProductRepository
	recentSearchesRepo   *repository.RecentSearchesRepository
}

func NewProductService(db *sqlx.DB, dbConfig config.DBConfig, orderConfig config.OrderCloudConfig, cacheClient *cache.Cache) *ProductService {
	return &ProductService{
		priceScheduleService: NewPriceScheduleService(),
		productRepo:          repository.NewProductRepository(db, dbConfig, orderConfig),
		categoryProductRepo:  repository.NewCategoryProductRepository(orderConfig),
		cacheClient:          cacheClient,
		recentSearchesRepo:   repository.NewRecentSearchesRepository(db, dbConfig),
	}
}

func (svc *ProductService) GetProducts(catalogID, categoryID, supplierID, userID, page, pageSize, sortBy, search *string, isFavorite *bool, extraFilters map[string]interface{}, accessToken string) (model.ProductResponse, error) {
	params := repository.ProductParams{
		CatalogID:    GetString(catalogID),
		CategoryID:   GetString(categoryID),
		SupplierID:   GetString(supplierID),
		Page:         GetString(page),
		PageSize:     GetString(pageSize),
		SortBy:       GetString(sortBy),
		ExtraFilters: extraFilters,
	}

	// get user favorite products from database
	userProductFavorites, err := svc.productRepo.GetFavoriteProducts(userID, page, pageSize)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// get favorite products from order cloud
	if userID != nil && isFavorite != nil && *isFavorite {
		favoriteProductIDs := getFavoriteProductIDs(userProductFavorites)
		// no favorite products found for user
		if favoriteProductIDs == "" {
			return model.ProductResponse{}, nil
		}

		if params.ExtraFilters == nil {
			params.ExtraFilters = make(map[string]interface{})
		}

		params.ExtraFilters["ID"] = favoriteProductIDs
	}

	// handle search
	if search != nil {
		params.Search = *search
		params.SearchOn = "ID,Name,Description"

		// insert search details to database
		_, err = svc.recentSearchesRepo.SaveRecentSearch(userID, search)
		if err != nil {
			return model.ProductResponse{}, err
		}
	}

	// get products from order cloud
	products, err := svc.productRepo.GetProducts(params, accessToken)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// add isFavorite field to products
	for i, product := range products.Items {
		for _, userProductFavorite := range userProductFavorites {
			if *product.ID == *userProductFavorite.ProductID {
				products.Items[i].IsFavorite = true
			}
		}
	}

	// prepare product ids
	productIDs := make([]string, len(products.Items))
	for i, product := range products.Items {
		productIDs[i] = GetString(product.ID)
	}

	// get price schedules
	priceSchedules, err := svc.priceScheduleService.GetPriceSchedules(productIDs, nil, nil, accessToken)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// add price schedules to products
	for i, product := range products.Items {
		for _, priceSchedule := range priceSchedules.Items {
			if GetString(product.ID) == GetString(priceSchedule.ID) {
				products.Items[i].PriceSchedule = priceSchedule
			}
		}
	}

	return products, nil
}

func (svc *ProductService) GetSimilarProducts(productID string, userID, page, pageSize *string, accessToken string) (model.ProductResponse, error) {
	categoryProduct, err := svc.categoryProductRepo.GetCategoryProductByProductID(productID, accessToken)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// get user favorite products from database
	userProductFavorites, err := svc.productRepo.GetFavoriteProducts(userID, page, pageSize)
	if err != nil {
		return model.ProductResponse{}, err
	}

	sameCategoryProducts, err := svc.productRepo.GetProducts(repository.ProductParams{
		CategoryID: categoryProduct.CategoryID,
		Page:       GetString(page),
		PageSize:   GetString(pageSize),
	}, accessToken)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// add isFavorite field to products
	for i, product := range sameCategoryProducts.Items {
		for _, userProductFavorite := range userProductFavorites {
			if *product.ID == *userProductFavorite.ProductID {
				sameCategoryProducts.Items[i].IsFavorite = true
			}
		}
	}

	// prepare product ids
	productIDs := make([]string, len(sameCategoryProducts.Items))
	for i, product := range sameCategoryProducts.Items {
		productIDs[i] = GetString(product.ID)
	}

	// get price schedules
	priceSchedules, err := svc.priceScheduleService.GetPriceSchedules(productIDs, nil, nil, accessToken)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// add price schedules to products
	for i, product := range sameCategoryProducts.Items {
		for _, priceSchedule := range priceSchedules.Items {
			if GetString(product.ID) == GetString(priceSchedule.ID) {
				sameCategoryProducts.Items[i].PriceSchedule = priceSchedule
			}
		}
	}

	return sameCategoryProducts, nil
}

func (svc *ProductService) GetRecommendProducts(productID string, page, pageSize *string, accessToken string) (model.ProductResponseV2, error) {
	categoryProduct, err := svc.categoryProductRepo.GetCategoryProductByProductID(productID, accessToken)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	sameCategoryProducts, err := svc.productRepo.GetProductsOrderCloudV2(repository.ProductParams{
		CategoryID: categoryProduct.CategoryID,
		Page:       GetString(page),
		PageSize:   GetString(pageSize),
	}, accessToken)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	// prepare product ids
	productIDs := make([]string, len(sameCategoryProducts.Items))
	for i, product := range sameCategoryProducts.Items {
		productIDs[i] = GetString(product.ID)
	}

	// get price schedules
	priceSchedules, err := svc.priceScheduleService.GetPriceSchedules(productIDs, nil, nil, accessToken)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	// add price schedules to products
	for i, product := range sameCategoryProducts.Items {
		for _, priceSchedule := range priceSchedules.Items {
			if GetString(product.ID) == GetString(priceSchedule.ID) {
				sameCategoryProducts.Items[i].PriceSchedule = &model.NewProductPriceSchedule{
					ID:                    priceSchedule.ID,
					OwnerID:               priceSchedule.OwnerID,
					Name:                  priceSchedule.Name,
					ApplyTax:              priceSchedule.ApplyTax,
					ApplyShipping:         priceSchedule.ApplyShipping,
					MinQuantity:           priceSchedule.MinQuantity,
					MaxQuantity:           priceSchedule.MaxQuantity,
					UseCumulativeQuantity: priceSchedule.UseCumulativeQuantity,
					RestrictedQuantity:    priceSchedule.RestrictedQuantity,
					Currency:              priceSchedule.Currency,
					SaleStart:             priceSchedule.SaleStart,
					SaleEnd:               priceSchedule.SaleEnd,
					IsOnSale:              priceSchedule.IsOnSale,
					PriceBreaks:           priceSchedule.PriceBreaks,
				}
			}
		}
	}

	return sameCategoryProducts, nil
}

func (svc *ProductService) GetProduct(productID string, userID *string, accessToken string) (model.ProductItem, error) {
	// get user favorite products from database
	userProductFavorites, err := svc.productRepo.GetFavoriteProducts(userID, nil, nil)
	if err != nil {
		return model.ProductItem{}, err
	}

	product, err := svc.productRepo.GetProduct(productID, accessToken)
	if err != nil {
		return model.ProductItem{}, err
	}

	// add isFavorite field to products
	for _, userProductFavorite := range userProductFavorites {
		if *product.ID == *userProductFavorite.ProductID {
			product.IsFavorite = true
		}
	}

	// get price schedules
	priceSchedules, err := svc.priceScheduleService.GetPriceSchedules([]string{productID}, nil, nil, accessToken)
	if err != nil {
		return model.ProductItem{}, err
	}

	// add price schedules to products
	for _, priceSchedule := range priceSchedules.Items {
		if GetString(product.ID) == GetString(priceSchedule.ID) {
			product.PriceSchedule = priceSchedule
		}
	}

	return product, nil
}

func (svc *ProductService) FavoriteProduct(userID *string, productID string, isFavorite bool) (*model.UserProductFavorite, error) {
	if isFavorite {
		userProductFavourite, err := svc.productRepo.SaveFavoriteProduct(userID, productID)
		if err != nil {
			return nil, err
		}

		return userProductFavourite, nil
	} else {
		err := svc.productRepo.DeleteFavoriteProduct(userID, productID)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (svc *ProductService) GetTrendingProducts(accessToken string) (model.ProductResponse, error) {
	cacheResp, err := svc.cacheClient.Get(TrendingProductCacheKey)
	if err != nil || cacheResp == nil {
		return model.ProductResponse{}, err
	}

	// unmarshal trending products to trending product struct
	var trendingProducts []model.TrendingProduct
	err = mapstructure.Decode(cacheResp, &trendingProducts)
	if err != nil {
		return model.ProductResponse{}, err
	}

	// get trending product ids
	trendingProductIDs := make([]string, len(trendingProducts))
	for i := range trendingProducts {
		trendingProductIDs[i] = *trendingProducts[i].ProductID
	}

	// prepare params for order cloud
	params := repository.ProductParams{
		ExtraFilters: map[string]interface{}{
			"ID": strings.Join(trendingProductIDs, "|"),
		},
	}

	// get products from order cloud
	products, err := svc.productRepo.GetProducts(params, accessToken)
	if err != nil {
		return model.ProductResponse{}, err
	}

	return products, nil
}

func (svc *ProductService) StartTrendingProductsProcessor(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()

	fmt.Println("starting trending products processor")

	// update trending products cache initially
	err := svc.updateTrendingProductCache()
	if err != nil {
		fmt.Println("error occurred while updating trending products cache initially:", err)
	}

	pullDuration := time.Duration(60) * time.Minute
	tick := time.NewTicker(pullDuration)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopped trending products processor")
			return
		case <-tick.C:
			fmt.Println("fetching trending products from database")

			// update trending products cache periodically
			err := svc.updateTrendingProductCache()
			if err != nil {
				continue
			}
		}
	}
}

func (svc *ProductService) updateTrendingProductCache() error {
	// get trending products from database
	trendingProducts, err := svc.productRepo.GetTrendingProducts(10)
	if err != nil {
		fmt.Println("error occurred while fetching trending products from database:", err)
		return err
	}

	fmt.Printf("fetched %d trending products from database\n", len(trendingProducts))

	// update trending products cache
	err = svc.cacheClient.Set(TrendingProductCacheKey, trendingProducts, 0)
	if err != nil {
		fmt.Println("error occurred while updating trending products cache:", err)
		return err
	}

	return nil
}

func getFavoriteProductIDs(userProductFavorites []model.UserProductFavorite) string {
	var productIDs []string
	for _, userProductFavorite := range userProductFavorites {
		if userProductFavorite.ProductID != nil {
			productIDs = append(productIDs, *userProductFavorite.ProductID)
		}
	}

	result := strings.Join(productIDs, "|")

	return result
}

func (svc *ProductService) GetProductsV2(catalogID, categoryID, supplierID, userID, page, pageSize, sortBy, search *string, isFavorite *bool, extraFilters map[string]interface{}, accessToken string) (model.ProductResponseV2, error) {
	params := repository.ProductParams{
		CatalogID:    GetString(catalogID),
		CategoryID:   GetString(categoryID),
		SupplierID:   GetString(supplierID),
		Page:         GetString(page),
		PageSize:     GetString(pageSize),
		SortBy:       GetString(sortBy),
		ExtraFilters: extraFilters,
	}

	// get user favorite products from database
	userProductFavorites, err := svc.productRepo.GetFavoriteProducts(userID, nil, nil)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	// get favorite products from order cloud
	if userID != nil && isFavorite != nil && *isFavorite {
		favoriteProductIDs := getFavoriteProductIDs(userProductFavorites)
		// no favorite products found for user
		if favoriteProductIDs == "" {
			return model.ProductResponseV2{}, nil
		}

		if params.ExtraFilters == nil {
			params.ExtraFilters = make(map[string]interface{})
		}

		params.ExtraFilters["ID"] = favoriteProductIDs
	}

	// handle search
	if search != nil {
		params.Search = *search
		params.SearchOn = "ID,Name,Description"

		// insert search details to database
		_, err = svc.recentSearchesRepo.SaveRecentSearch(userID, search)
		if err != nil {
			return model.ProductResponseV2{}, err
		}
	}

	// get products from order cloud
	products, err := svc.productRepo.GetProductsV2(params, accessToken)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	// add isFavorite field to products
	for i, product := range products.Items {
		for _, userProductFavorite := range userProductFavorites {
			if *product.Product.ID == *userProductFavorite.ProductID {
				products.Items[i].Product.IsFavorite = true
			}
		}
	}

	return products, nil
}

func (svc *ProductService) GetProductV2(productID string, userID *string, accessToken string) (model.LatestProductItems, error) {
	// get user favorite products from database
	userProductFavorites, err := svc.productRepo.GetFavoriteProducts(userID, nil, nil)
	if err != nil {
		return model.LatestProductItems{}, err
	}

	product, err := svc.productRepo.GetProductV2(productID, accessToken)
	if err != nil {
		return model.LatestProductItems{}, err
	}

	// add isFavorite field to products
	for _, userProductFavorite := range userProductFavorites {
		if *product.Product.ID == *userProductFavorite.ProductID {
			product.Product.IsFavorite = true
		}
	}

	return product, nil
}

func (svc *ProductService) GetProductFilterMiddleWare(searchOn string, accessToken string) ([]*model.ProductFilter, error) {
	// get products filters from .net middleware
	productFilter, err := svc.productRepo.FetchProductFilters(searchOn, accessToken)
	if err != nil {
		return []*model.ProductFilter{}, err
	}
	return productFilter, nil
}
