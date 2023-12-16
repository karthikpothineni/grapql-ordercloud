package service

import (
	"errors"
	"mpmy-product-service/graph/model"
	"mpmy-product-service/repository"
	"testing"

	"github.com/google/uuid"
)

func TestGetRecommendProducts(t *testing.T) {

	//data test

	var accessToken = uuid.New().String()
	var productID = uuid.New().String()
	var categoryID = uuid.New().String()
	var productIDs = []string{}

	var page *string
	var pageSize *string

	var latestProductItems []*model.LatestProductItems
	var priceScheduleItem []*model.PriceScheduleItem

	//create list product id for GetPriceSchedules service
	//create list product item for GetProductsOrderCloudV2 repo
	//create list schedule item of product id for GetPriceSchedules service
	for i := 0; i < 3; i++ {
		var id = uuid.New().String()
		productIDs = append(productIDs, id)
		latestProductItems = append(latestProductItems, &model.LatestProductItems{
			ID: &id,
		})
		priceScheduleItem = append(priceScheduleItem, &model.PriceScheduleItem{
			ID: &id,
		})
	}

	var categoryProductRepositoryMock = &repository.CategoryProductRepositoryMock{}
	var productRepositoryMock = &repository.ProductRepositoryMock{}
	var priceScheduleServiceMock = &PriceScheduleServiceMock{}

	categoryProductRepositoryMock.On("GetCategoryProductByProductID", productID, accessToken).Return(repository.CategoryProductItem{
		ProductID:  productID,
		CategoryID: categoryID,
	}, nil)

	productRepositoryMock.On("GetProductsOrderCloudV2", repository.ProductParams{CategoryID: categoryID, Page: "", PageSize: ""}, accessToken).Return(model.ProductResponseV2{
		Items: latestProductItems,
	}, nil)

	priceScheduleServiceMock.On("GetPriceSchedules", productIDs, page, pageSize, accessToken).Return(model.PriceScheduleResponse{
		Items: priceScheduleItem,
	}, nil)

	productService := &ProductService{
		productRepo:          productRepositoryMock,
		categoryProductRepo:  categoryProductRepositoryMock,
		priceScheduleService: priceScheduleServiceMock,
	}

	data, err := productService.GetRecommendProducts(productID, page, pageSize, accessToken)
	if err != nil {
		t.Error("test failed error: ", err)
	} else if len(data.Items) <= 0 {
		t.Error("test failed: items must return data")
	} else {
		//check result data valid
		for _, d := range data.Items {
			valid := false
			for _, pid := range productIDs {

				if *d.ID == pid {
					valid = true
					break
				}
			}

			if !valid {
				t.Error("test failed: return data not match")
			}
		}
	}
}

func TestGetRecommendProductsWithErrorCategoryRepo(t *testing.T) {

	//data test
	var errorTest = errors.New("new error")
	var accessToken = uuid.New().String()
	var productID = uuid.New().String()
	var categoryID = uuid.New().String()

	var page *string
	var pageSize *string

	var categoryProductRepositoryMock = &repository.CategoryProductRepositoryMock{}

	categoryProductRepositoryMock.On("GetCategoryProductByProductID", productID, accessToken).Return(repository.CategoryProductItem{
		ProductID:  productID,
		CategoryID: categoryID,
	}, errorTest)

	productService := &ProductService{
		categoryProductRepo: categoryProductRepositoryMock,
	}

	_, err := productService.GetRecommendProducts(productID, page, pageSize, accessToken)
	if err == nil || err != errorTest {
		t.Error("test failed error")
	}
}

func TestGetRecommendProductsWithErrorProductRepo(t *testing.T) {

	//data test
	var errorTest = errors.New("new error")
	var accessToken = uuid.New().String()
	var productID = uuid.New().String()
	var categoryID = uuid.New().String()

	var page *string
	var pageSize *string

	var latestProductItems []*model.LatestProductItems

	var categoryProductRepositoryMock = &repository.CategoryProductRepositoryMock{}
	var productRepositoryMock = &repository.ProductRepositoryMock{}

	categoryProductRepositoryMock.On("GetCategoryProductByProductID", productID, accessToken).Return(repository.CategoryProductItem{
		ProductID:  productID,
		CategoryID: categoryID,
	}, nil)

	productRepositoryMock.On("GetProductsOrderCloudV2", repository.ProductParams{CategoryID: categoryID, Page: "", PageSize: ""}, accessToken).Return(model.ProductResponseV2{
		Items: latestProductItems,
	}, errorTest)

	productService := &ProductService{
		productRepo:         productRepositoryMock,
		categoryProductRepo: categoryProductRepositoryMock,
	}

	_, err := productService.GetRecommendProducts(productID, page, pageSize, accessToken)
	if err == nil || err != errorTest {
		t.Error("test failed error")
	}
}

func TestGetRecommendProductsWithErrorPrinceSchedule(t *testing.T) {

	//data test
	var errorTest = errors.New("new error")
	var accessToken = uuid.New().String()
	var productID = uuid.New().String()
	var categoryID = uuid.New().String()
	var productIDs = []string{}

	var page *string
	var pageSize *string

	var latestProductItems []*model.LatestProductItems
	var priceScheduleItem []*model.PriceScheduleItem

	//create list product id for GetPriceSchedules service
	//create list product item for GetProductsOrderCloudV2 repo
	//create list schedule item of product id for GetPriceSchedules service
	for i := 0; i < 3; i++ {
		var id = uuid.New().String()
		productIDs = append(productIDs, id)
		latestProductItems = append(latestProductItems, &model.LatestProductItems{
			ID: &id,
		})
		priceScheduleItem = append(priceScheduleItem, &model.PriceScheduleItem{
			ID: &id,
		})
	}

	var categoryProductRepositoryMock = &repository.CategoryProductRepositoryMock{}
	var productRepositoryMock = &repository.ProductRepositoryMock{}
	var priceScheduleServiceMock = &PriceScheduleServiceMock{}

	categoryProductRepositoryMock.On("GetCategoryProductByProductID", productID, accessToken).Return(repository.CategoryProductItem{
		ProductID:  productID,
		CategoryID: categoryID,
	}, nil)

	productRepositoryMock.On("GetProductsOrderCloudV2", repository.ProductParams{CategoryID: categoryID, Page: "", PageSize: ""}, accessToken).Return(model.ProductResponseV2{
		Items: latestProductItems,
	}, nil)

	priceScheduleServiceMock.On("GetPriceSchedules", productIDs, page, pageSize, accessToken).Return(model.PriceScheduleResponse{
		Items: priceScheduleItem,
	}, errorTest)

	productService := &ProductService{
		productRepo:          productRepositoryMock,
		categoryProductRepo:  categoryProductRepositoryMock,
		priceScheduleService: priceScheduleServiceMock,
	}

	_, err := productService.GetRecommendProducts(productID, page, pageSize, accessToken)
	if err == nil || err != errorTest {
		t.Error("test failed error")
	}
}
