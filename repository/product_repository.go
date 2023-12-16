package repository

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"mpmy-product-service/config"
	"mpmy-product-service/graph/model"
	"mpmy-product-service/httprequest"
)

const (
	// FavoriteProductTableName is the name of favorite product table
	FavoriteProductTableName = "user_product_favorites"

	// TrendingProductTableName is the name of trending product table
	TrendingProductTableName = "trending_products"

	// DefaultCatalogID is the catalog id in order cloud for MY
	DefaultCatalogID = "zp-my"
)

type ProductParams struct {
	CatalogID    string                 `json:"catalogID"`
	CategoryID   string                 `json:"categoryID"`
	SupplierID   string                 `json:"supplierID"`
	Search       string                 `json:"search"`
	SearchOn     string                 `json:"searchOn"`
	SortBy       string                 `json:"sortBy"`
	ExtraFilters map[string]interface{} `json:"extraFilters"`
	Page         string                 `json:"page"`
	PageSize     string                 `json:"pageSize"`
}

type IProductRepository interface {
	GetProducts(params ProductParams, accessToken string) (model.ProductResponse, error)
	GetProduct(productID string, accessToken string) (model.ProductItem, error)
	SaveFavoriteProduct(userID *string, productID string) (*model.UserProductFavorite, error)
	DeleteFavoriteProduct(userID *string, productID string) error
	GetFavoriteProducts(userID, page, pageSize *string) ([]model.UserProductFavorite, error)
	GetProductsOrderCloudV2(params ProductParams, accessToken string) (model.ProductResponseV2, error)
	GetProductsV2(params ProductParams, accessToken string) (model.ProductResponseV2, error)
	GetProductV2(productID string, accessToken string) (model.LatestProductItems, error)
	GetTrendingProducts(limit int) ([]model.TrendingProduct, error)
	FetchProductFilters(search string, accessToken string) ([]*model.ProductFilter, error)
}

type ProductRepository struct {
	db                 *sqlx.DB
	dbConfig           config.DBConfig
	orderCloud         config.OrderCloudConfig
	httpRequestHandler *httprequest.RequestHandler
}

func NewProductRepository(db *sqlx.DB, dbConfig config.DBConfig, orderConfig config.OrderCloudConfig) *ProductRepository {
	return &ProductRepository{
		db:                 db,
		dbConfig:           dbConfig,
		orderCloud:         orderConfig,
		httpRequestHandler: httprequest.NewRequestHandler("ProductRepository"),
	}
}

func (repo *ProductRepository) GetProducts(params ProductParams, accessToken string) (model.ProductResponse, error) {
	if params.CatalogID == "" {
		params.CatalogID = DefaultCatalogID
	}

	// prepare request specifications
	requestURL := fmt.Sprintf("%s/%s", repo.orderCloud.OrderCloudEngine, "v1/products")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        requestURL,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"catalogID":  params.CatalogID,
			"categoryID": params.CategoryID,
			"supplierID": params.SupplierID,
			"search":     params.Search,
			"searchOn":   params.SearchOn,
		},
	}

	if params.Page != "" {
		requestSpecifications.Params["page"] = params.Page
	}

	if params.PageSize != "" {
		requestSpecifications.Params["pageSize"] = params.PageSize
	}

	if params.SortBy != "" {
		requestSpecifications.Params["sortBy"] = params.SortBy
	}

	// add extra filters from params
	if params.ExtraFilters == nil {
		params.ExtraFilters = map[string]interface{}{
			"Active": true,
		}
	} else {
		params.ExtraFilters["Active"] = true
	}
	if len(params.ExtraFilters) > 0 {
		for key, value := range params.ExtraFilters {
			requestSpecifications.Params[key] = value
		}
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.ProductResponse{}, fmt.Errorf("failed to fetch products")
	}

	var productResp model.ProductResponse

	err := json.Unmarshal(response, &productResp)
	if err != nil {
		return model.ProductResponse{}, err
	}

	return productResp, nil
}

func (repo *ProductRepository) GetProduct(productID string, accessToken string) (model.ProductItem, error) {
	url := fmt.Sprintf("%s/%s/%s", repo.orderCloud.OrderCloudEngine, "v1/products", productID)
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.ProductItem{}, fmt.Errorf("failed to fetch product")
	}
	var product model.ProductItem

	err := json.Unmarshal(response, &product)
	if err != nil {
		return model.ProductItem{}, err
	}

	return product, nil
}

func (repo *ProductRepository) SaveFavoriteProduct(userID *string, productID string) (*model.UserProductFavorite, error) {
	currentTime := time.Now().UTC()
	userProductFavourite := &model.UserProductFavorite{
		UserID:    userID,
		ProductID: &productID,
		CreatedAt: &currentTime,
	}

	query := "INSERT INTO " + repo.dbConfig.Schema + "." + FavoriteProductTableName + "(" +
		"user_id, product_id, created_at)" +
		"VALUES(:user_id, :product_id, :created_at) " +
		"RETURNING id"

	rows, err := repo.db.NamedQuery(query, userProductFavourite)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var lastInsertID int
	if rows.Next() {
		err = rows.Scan(&lastInsertID)
		if err != nil {
			return nil, err
		}
	}

	userProductFavourite.ID = &lastInsertID

	return userProductFavourite, nil
}

func (repo *ProductRepository) DeleteFavoriteProduct(userID *string, productID string) error {
	query := "DELETE FROM " + repo.dbConfig.Schema + "." + FavoriteProductTableName +
		" WHERE user_id = $1 AND product_id = $2"

	_, err := repo.db.Exec(query, userID, productID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductRepository) GetFavoriteProducts(userID, page, pageSize *string) ([]model.UserProductFavorite, error) {
	var userProductFavorites []model.UserProductFavorite
	var err error

	query := "SELECT * FROM " + repo.dbConfig.Schema + "." + FavoriteProductTableName +
		" WHERE user_id = $1"

	if page != nil && pageSize != nil {
		query += " LIMIT $2 OFFSET $3"

		limit := GetIntFromStringPointer(pageSize)
		offset := (GetIntFromStringPointer(page) * limit) - limit

		err = repo.db.Select(&userProductFavorites, query, userID, limit, offset)
	} else {
		err = repo.db.Select(&userProductFavorites, query, userID)
	}

	if err != nil {
		return nil, err
	}

	return userProductFavorites, nil
}

func (repo *ProductRepository) GetProductsOrderCloudV2(params ProductParams, accessToken string) (model.ProductResponseV2, error) {
	if params.CatalogID == "" {
		params.CatalogID = DefaultCatalogID
	}

	requestURL := fmt.Sprintf("%s/%s", repo.orderCloud.OrderCloudEngine, "v1/products")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        requestURL,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"catalogID":  params.CatalogID,
			"categoryID": params.CategoryID,
			"supplierID": params.SupplierID,
			"search":     params.Search,
			"searchOn":   params.SearchOn,
		},
	}

	if params.Page != "" {
		requestSpecifications.Params["page"] = params.Page
	}

	if params.PageSize != "" {
		requestSpecifications.Params["pageSize"] = params.PageSize
	}

	if params.SortBy != "" {
		requestSpecifications.Params["sortBy"] = params.SortBy
	}

	// add extra filters from params
	if params.ExtraFilters == nil {
		params.ExtraFilters = map[string]interface{}{
			"Active": true,
		}
	} else {
		params.ExtraFilters["Active"] = true
	}
	if len(params.ExtraFilters) > 0 {
		for key, value := range params.ExtraFilters {
			requestSpecifications.Params[key] = value
		}
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.ProductResponseV2{}, fmt.Errorf("failed to fetch products")
	}

	var productResp model.ProductResponseV2

	err := json.Unmarshal(response, &productResp)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	return productResp, nil
}

func (repo *ProductRepository) GetProductsV2(params ProductParams, accessToken string) (model.ProductResponseV2, error) {
	if params.CatalogID == "" {
		params.CatalogID = DefaultCatalogID
	}

	requestURL := fmt.Sprintf("%s/%s", repo.orderCloud.SellerCenterMiddleware, "products")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        requestURL,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"catalogID":  params.CatalogID,
			"categoryID": params.CategoryID,
			"supplierID": params.SupplierID,
			"search":     params.Search,
			"searchOn":   params.SearchOn,
		},
	}

	if params.Page != "" {
		requestSpecifications.Params["page"] = params.Page
	}

	if params.PageSize != "" {
		requestSpecifications.Params["pageSize"] = params.PageSize
	}

	if params.SortBy != "" {
		requestSpecifications.Params["sortBy"] = params.SortBy
	}

	// add extra filters from params
	if params.ExtraFilters == nil {
		params.ExtraFilters = map[string]interface{}{
			"Active": true,
		}
	} else {
		params.ExtraFilters["Active"] = true
	}
	if len(params.ExtraFilters) > 0 {
		for key, value := range params.ExtraFilters {
			requestSpecifications.Params[key] = value
		}
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.ProductResponseV2{}, fmt.Errorf("failed to fetch products")
	}

	var productResp model.ProductResponseV2

	err := json.Unmarshal(response, &productResp)
	if err != nil {
		return model.ProductResponseV2{}, err
	}

	return productResp, nil
}

func (repo *ProductRepository) GetProductV2(productID string, accessToken string) (model.LatestProductItems, error) {
	url := fmt.Sprintf("%s/%s/%s", repo.orderCloud.SellerCenterMiddleware, "products", productID)
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.LatestProductItems{}, fmt.Errorf("failed to fetch product")
	}
	var product model.LatestProductItems

	err := json.Unmarshal(response, &product)
	if err != nil {
		return model.LatestProductItems{}, err
	}
	return product, nil
}

func (repo *ProductRepository) GetTrendingProducts(limit int) ([]model.TrendingProduct, error) {
	var trendingProducts []model.TrendingProduct
	var err error

	createdAtStart := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")
	createdAtEnd := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")

	query := "SELECT product_id, COUNT(product_id) AS order_count, sum(quantity) as quantity FROM " + repo.dbConfig.Schema + "." + TrendingProductTableName +
		" WHERE created_at between $1 AND $2 GROUP BY product_id ORDER BY order_count DESC LIMIT $3"

	err = repo.db.Select(&trendingProducts, query, createdAtStart, createdAtEnd, limit)
	if err != nil {
		return nil, err
	}

	return trendingProducts, nil
}

// /products/Product_filter?Search=CountryOfOrigin
func (repo *ProductRepository) FetchProductFilters(search string, accessToken string) ([]*model.ProductFilter, error) {
	url := fmt.Sprintf("%s/%s", repo.orderCloud.SellerCenterMiddleware, "products/Product_filter")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"Search": search,
		},
	}
	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return []*model.ProductFilter{}, fmt.Errorf("failed to fetch product")
	}
	var product []*model.ProductFilter

	err := json.Unmarshal(response, &product)
	if err != nil {
		return []*model.ProductFilter{}, err
	}
	return product, nil
}
