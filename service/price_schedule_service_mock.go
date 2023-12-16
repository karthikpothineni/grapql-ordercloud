package service

import (
	"mpmy-product-service/graph/model"
	"mpmy-product-service/repository"

	"github.com/stretchr/testify/mock"
)

type PriceScheduleServiceMock struct {
	priceScheduleRepo *repository.PriceScheduleRepository
	mock.Mock
}

func (svc *PriceScheduleServiceMock) GetPriceSchedule(productID string, page, pageSize *string, accessToken string) (model.PriceScheduleResponse, error) {
	params := repository.PriceScheduleParams{
		Search:   productID,
		SearchOn: "ID",
		Page:     GetString(page),
		PageSize: GetString(pageSize),
	}

	priceSchedules, err := svc.priceScheduleRepo.GetPriceSchedules(params, accessToken)
	if err != nil {
		return model.PriceScheduleResponse{}, err
	}

	return priceSchedules, nil
}

func (svc *PriceScheduleServiceMock) GetPriceSchedules(productIDs []string, page, pageSize *string, accessToken string) (model.PriceScheduleResponse, error) {

	args := svc.Called(productIDs, page, pageSize, accessToken)

	return args.Get(0).(model.PriceScheduleResponse), args.Error(1)
}
