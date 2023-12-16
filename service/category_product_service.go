package service

import (
	"mpmy-product-service/config"
	"mpmy-product-service/repository"
)

type CategoryProductService struct {
	categoryProductRepo *repository.CategoryProductRepository
}

func NewCategoryProductService(orderConfig config.OrderCloudConfig) *CategoryProductService {
	return &CategoryProductService{
		categoryProductRepo: repository.NewCategoryProductRepository(orderConfig),
	}
}

func (svc *CategoryProductService) GetCategoryProducts(params repository.CategoryProductParams, accessToken string) (repository.CategoryProductResponse, error) {
	categoryProducts, err := svc.categoryProductRepo.GetCategoryProducts(params, accessToken)
	if err != nil {
		return repository.CategoryProductResponse{}, err
	}

	return categoryProducts, nil
}
