package service

import (
	"mpmy-product-service/config"
	"mpmy-product-service/constants"
	"mpmy-product-service/graph/model"
	"mpmy-product-service/repository"
)

type CategoryService struct {
	CategoryRepo *repository.CategoryRepository
}

func NewCategoryService(orderConfig config.OrderCloudConfig) *CategoryService {
	return &CategoryService{
		CategoryRepo: repository.NewCategoryRepository(orderConfig),
	}
}

func (svc *CategoryService) GetCategories(catalog, depth *string, accessToken string) (model.CategoryResponse, error) {
	catalogID := GetString(catalog)
	depthVal := GetString(depth)
	if catalogID == "" {
		catalogID = constants.CatalogID
	}
	if depthVal == "" {
		depthVal = constants.Depth
	}
	params := repository.GetCategoryParams{
		CatalogID: catalogID,
		Depth:     depthVal,
	}
	categoryProducts, err := svc.CategoryRepo.FetchCategories(params, accessToken)
	if err != nil {
		return model.CategoryResponse{}, err
	}

	return categoryProducts, nil
}
