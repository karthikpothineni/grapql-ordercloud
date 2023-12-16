package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mpmy-product-service/config"
	"mpmy-product-service/httprequest"
)

const (
	MYCatalogID = "zp-my"
)

type CategoryProductParams struct {
	CategoryID string `json:"categoryID"`
	ProductID  string `json:"productID"`
	Page       string `json:"page"`
	PageSize   string `json:"pageSize"`
}

type CategoryProductResponse struct {
	Meta  categoryProductMeta   `json:"Meta"`
	Items []CategoryProductItem `json:"Items"`
}

type categoryProductMeta struct {
	Page        int    `json:"Page"`
	PageSize    int    `json:"PageSize"`
	TotalCount  int    `json:"TotalCount"`
	TotalPages  int    `json:"TotalPages"`
	ItemRange   []int  `json:"ItemRange"`
	NextPageKey string `json:"NextPageKey"`
}

type CategoryProductItem struct {
	CategoryID string `json:"CategoryID"`
	ProductID  string `json:"ProductID"`
	ListOrder  int    `json:"ListOrder"`
}

type ICategoryProductRepository interface {
	GetCategoryProducts(params CategoryProductParams, accessToken string) (CategoryProductResponse, error)
	GetCategoryProductByProductID(productID string, accessToken string) (CategoryProductItem, error)
}

type CategoryProductRepository struct {
	orderCloud         config.OrderCloudConfig
	httpRequestHandler *httprequest.RequestHandler
}

func NewCategoryProductRepository(orderConfig config.OrderCloudConfig) *CategoryProductRepository {
	return &CategoryProductRepository{
		orderCloud:         orderConfig,
		httpRequestHandler: httprequest.NewRequestHandler("CategoryProductRepository"),
	}
}

func (repo *CategoryProductRepository) GetCategoryProducts(params CategoryProductParams, accessToken string) (CategoryProductResponse, error) {
	// prepare request specifications
	url := fmt.Sprintf("%s/%s", repo.orderCloud.OrderCloudEngine, "v1/catalogs/"+MYCatalogID+"/categories/productassignments")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"categoryID": params.CategoryID,
			"productID":  params.ProductID,
			"page":       params.Page,
			"pageSize":   params.PageSize,
		},
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return CategoryProductResponse{}, fmt.Errorf("failed to fetch category product assignments")
	}

	var assignmentResp CategoryProductResponse

	err := json.Unmarshal(response, &assignmentResp)
	if err != nil {
		return CategoryProductResponse{}, err
	}

	return assignmentResp, nil
}

func (repo *CategoryProductRepository) GetCategoryProductByProductID(productID string, accessToken string) (CategoryProductItem, error) {
	// prepare request specifications
	url := fmt.Sprintf("%s/%s", repo.orderCloud.OrderCloudEngine, "v1/catalogs/"+MYCatalogID+"/categories/productassignments")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"productID": productID,
		},
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return CategoryProductItem{}, fmt.Errorf("failed to fetch category product assignment")
	}

	var assignmentResp CategoryProductResponse

	err := json.Unmarshal(response, &assignmentResp)
	if err != nil {
		return CategoryProductItem{}, err
	}

	if len(assignmentResp.Items) == 0 {
		return CategoryProductItem{}, fmt.Errorf("no category product assignment found")
	}

	return assignmentResp.Items[0], nil
}
