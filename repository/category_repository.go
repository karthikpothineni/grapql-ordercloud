package repository

import (
	"encoding/json"
	"fmt"
	"mpmy-product-service/config"
	"mpmy-product-service/constants"
	"mpmy-product-service/graph/model"
	"mpmy-product-service/httprequest"
	"net/http"
)

type GetCategoryParams struct {
	Depth     string `json:"depth"`
	CatalogID string `json:"catalogID"`
}
type CategoryRepository struct {
	orderCloud         config.OrderCloudConfig
	httpRequestHandler *httprequest.RequestHandler
}

func NewCategoryRepository(orderConfig config.OrderCloudConfig) *CategoryRepository {
	return &CategoryRepository{
		orderCloud:         orderConfig,
		httpRequestHandler: httprequest.NewRequestHandler("CategoryRepository"),
	}
}

func (repo *CategoryRepository) FetchCategories(params GetCategoryParams, accessToken string) (model.CategoryResponse, error) {
	// prepare request specifications
	url := fmt.Sprintf("%s/%s", repo.orderCloud.OrderCloudEngine, "v1/catalogs/"+params.CatalogID+"/categories")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"depth":    params.Depth,
			"pageSize": constants.MinimumCategoriesShown,
		},
	}
	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.CategoryResponse{}, fmt.Errorf("%s", response)
	}
	var assignmentResp model.CategoryResponse
	err := json.Unmarshal(response, &assignmentResp)
	if err != nil {
		return model.CategoryResponse{}, err
	}
	depthTwoCategory := make(map[string]*model.CategoryItems)
	childCategory := make(map[string]*model.CategoryItems)
	parentCategory := make(map[string]*model.CategoryItems)
	for _, category := range assignmentResp.Items {
		if category.ParentID == "" {
			parentCategory[category.ID] = category
			continue
		} else if category.ChildCount == 0 && category.ParentID != "" {
			childCategory[category.ID] = category
			continue
		}
		depthTwoCategory[category.ID] = category
	}
	depthTwoCategory, parentCategory = reOrderingDepthCategory(depthTwoCategory, childCategory, parentCategory)
	parentCategory = parentCategoryResponse(depthTwoCategory, parentCategory)
	var responseRestructure []*model.CategoryItems
	for _, item := range parentCategory {
		responseRestructure = append(responseRestructure, item)
	}
	assignmentResp.Items = responseRestructure
	return assignmentResp, nil
}

func reOrderingDepthCategory(depthTwoCategory, childCategory, parentCategory map[string]*model.CategoryItems) (map[string]*model.CategoryItems, map[string]*model.CategoryItems) {
	for _, val := range childCategory {
		if depthTwoCategorydata, ok := depthTwoCategory[val.ParentID]; ok {
			depthTwoCategorydata.ChildData = append(depthTwoCategorydata.ChildData, val)
			depthTwoCategory[val.ParentID] = depthTwoCategorydata
			defer delete(childCategory, val.ID)
			continue
		} else if parentCategorydata, ok := parentCategory[val.ParentID]; ok {
			parentCategorydata.ChildData = append(parentCategorydata.ChildData, val)
			parentCategory[val.ParentID] = parentCategorydata
			defer delete(childCategory, val.ID)
			continue
		}
	}
	return depthTwoCategory, parentCategory
}
func parentCategoryResponse(depthTwoCategory, parentCategory map[string]*model.CategoryItems) map[string]*model.CategoryItems {
	for key, val := range depthTwoCategory {
		if parentCategorydata, ok := parentCategory[val.ParentID]; ok {
			parentCategorydata.ChildData = append(parentCategorydata.ChildData, val)
			parentCategory[val.ParentID] = parentCategorydata
			defer delete(depthTwoCategory, key)
		}
	}
	return parentCategory
}
