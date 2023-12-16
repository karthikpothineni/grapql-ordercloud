package service

import (
	"context"
	"fmt"
	"mpmy-product-service/config"

	"mpmy-product-service/repository"
	"time"
)

type LoginService struct {
	AccessToken string
	loginRepo   *repository.LoginRepository
}

func NewLoginService() *LoginService {
	return &LoginService{
		loginRepo: repository.NewLoginRepository(),
	}
}

func (svc *LoginService) StartAccessTokenFetcher(ctx context.Context, orderCloudConfig config.OrderCloudConfig) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()

	fmt.Println("starting access token fetcher")

	// fetch access token before the ticker starts
	accessToken, err := svc.loginRepo.GetAccessToken(orderCloudConfig)
	if err != nil {
		fmt.Println("failed to fetch access token, error: ", err)
	} else {
		svc.AccessToken = accessToken
		fmt.Println("access token fetched")
	}

	pullDuration := time.Duration(orderCloudConfig.AccessTokenFetchDuration) * time.Minute
	tick := time.NewTicker(pullDuration)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopped access token fetcher")
			return
		case <-tick.C:
			fmt.Println("fetching access token")
			accessToken, err = svc.loginRepo.GetAccessToken(orderCloudConfig)
			if err != nil {
				fmt.Println("failed to fetch access token")
			} else {
				svc.AccessToken = accessToken
				fmt.Println("access token fetched")
			}
		}
	}
}
