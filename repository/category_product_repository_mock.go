package repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mpmy-product-service/config"
	"mpmy-product-service/httprequest"

	"github.com/stretchr/testify/mock"
)

var CategoryProductRepositoryMockData = mock.Mock{}

type CategoryProductRepositoryMock struct {
	orderCloud         config.OrderCloudConfig
	httpRequestHandler *httprequest.RequestHandler
	mock.Mock
}

func (repo *CategoryProductRepositoryMock) GetCategoryProducts(params CategoryProductParams, accessToken string) (CategoryProductResponse, error) {
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

func (repo *CategoryProductRepositoryMock) GetCategoryProductByProductID(productID string, accessToken string) (CategoryProductItem, error) {

	args := repo.Called(productID, accessToken)

	return args.Get(0).(CategoryProductItem), args.Error(1)
}
