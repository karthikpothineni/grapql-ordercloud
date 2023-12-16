package service

import (
	"mpmy-product-service/graph/model"
	"mpmy-product-service/repository"
	"strings"
)

type IPriceScheduleService interface {
	GetPriceSchedule(productID string, page, pageSize *string, accessToken string) (model.PriceScheduleResponse, error)
	GetPriceSchedules(productIDs []string, page, pageSize *string, accessToken string) (model.PriceScheduleResponse, error)
}

type PriceScheduleService struct {
	priceScheduleRepo *repository.PriceScheduleRepository
}

func NewPriceScheduleService() *PriceScheduleService {
	return &PriceScheduleService{
		priceScheduleRepo: repository.NewPriceScheduleRepository(),
	}
}

func (svc *PriceScheduleService) GetPriceSchedule(productID string, page, pageSize *string, accessToken string) (model.PriceScheduleResponse, error) {
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

func (svc *PriceScheduleService) GetPriceSchedules(productIDs []string, page, pageSize *string, accessToken string) (model.PriceScheduleResponse, error) {
	params := repository.PriceScheduleParams{
		ExtraFilters: map[string]interface{}{
			"ID": strings.Join(productIDs, "|"),
		},
		Page:     GetString(page),
		PageSize: GetString(pageSize),
	}

	priceSchedules, err := svc.priceScheduleRepo.GetPriceSchedules(params, accessToken)
	if err != nil {
		return model.PriceScheduleResponse{}, err
	}

	return priceSchedules, nil
}
