package repository

import (
	"encoding/json"
	"fmt"
	"mpmy-product-service/graph/model"
	"net/http"

	"mpmy-product-service/constants"
	"mpmy-product-service/httprequest"
)

type PriceScheduleParams struct {
	Search       string                 `json:"search"`
	SearchOn     string                 `json:"searchOn"`
	ExtraFilters map[string]interface{} `json:"extraFilters"`
	Page         string                 `json:"page"`
	PageSize     string                 `json:"pageSize"`
}

type IPriceScheduleRepository interface {
	GetPriceSchedules(params PriceScheduleParams, accessToken string) (model.PriceScheduleResponse, error)
}

type PriceScheduleRepository struct {
	httpRequestHandler *httprequest.RequestHandler
}

func NewPriceScheduleRepository() *PriceScheduleRepository {
	return &PriceScheduleRepository{
		httpRequestHandler: httprequest.NewRequestHandler("PriceScheduleRepository"),
	}
}

func (repo *PriceScheduleRepository) GetPriceSchedules(params PriceScheduleParams, accessToken string) (model.PriceScheduleResponse, error) {
	// prepare request specifications
	url := fmt.Sprintf("%s/%s", constants.OrderCloudEngine, "v1/priceschedules")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod: http.MethodGet,
		URL:        url,
		Headers:    map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)},
		Params: map[string]interface{}{
			"search":   params.Search,
			"searchOn": params.SearchOn,
		},
	}

	if params.Page != "" {
		requestSpecifications.Params["page"] = params.Page
	}

	if params.PageSize != "" {
		requestSpecifications.Params["pageSize"] = params.PageSize
	}

	// add extra filters
	if params.ExtraFilters != nil && len(params.ExtraFilters) > 0 {
		for key, value := range params.ExtraFilters {
			requestSpecifications.Params[key] = value
		}
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return model.PriceScheduleResponse{}, fmt.Errorf("failed to fetch price schedules")
	}

	var assignmentResp model.PriceScheduleResponse

	err := json.Unmarshal(response, &assignmentResp)
	if err != nil {
		return model.PriceScheduleResponse{}, err
	}

	return assignmentResp, nil
}
