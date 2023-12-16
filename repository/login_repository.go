package repository

import (
	"encoding/json"
	"fmt"
	"mpmy-product-service/config"
	"net/http"

	"mpmy-product-service/constants"
	"mpmy-product-service/httprequest"
)

const (
	GrantTypePassword = "password"
)

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginRepository struct {
	httpRequestHandler *httprequest.RequestHandler
}

func NewLoginRepository() *LoginRepository {
	return &LoginRepository{
		httpRequestHandler: httprequest.NewRequestHandler("LoginRepository"),
	}
}

func (repo *LoginRepository) GetAccessToken(orderCloudConfig config.OrderCloudConfig) (string, error) {
	// prepare request specifications
	url := fmt.Sprintf("%s/%s", constants.OrderCloudEngine, "oauth/token")
	requestSpecifications := &httprequest.RequestSpecifications{
		HTTPMethod:  http.MethodPost,
		URL:         url,
		RequestType: "form",
		Params: map[string]interface{}{
			"client_id":     orderCloudConfig.ClientID,
			"client_secret": orderCloudConfig.ClientSecret,
			"grant_type":    GrantTypePassword,
			"username":      orderCloudConfig.Username,
			"password":      orderCloudConfig.Password,
		},
	}

	// make request
	statusCode, response, _ := repo.httpRequestHandler.MakeRequest(requestSpecifications)
	if statusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch access token")
	}

	var loginResp loginResponse

	err := json.Unmarshal(response, &loginResp)
	if err != nil {
		return "", err
	}

	return loginResp.AccessToken, nil
}
