package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	General    GeneralConfig
	OrderCloud OrderCloudConfig
	DB         DBConfig
}

type GeneralConfig struct {
	Name           string
	Port           string
	StaticToken    string
	ByPassSecurity bool
}

type OrderCloudConfig struct {
	ClientID                 string
	ClientSecret             string
	Username                 string
	Password                 string
	AccessTokenFetchDuration int
	OrderCloudEngine         string
	SellerCenterMiddleware   string
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Schema   string
	SSLMode  string
}

// Init - prepares the config from environmental variables
func Init() AppConfig {
	appConfig := AppConfig{
		General: GeneralConfig{
			Name: panicIfEmpty(os.Getenv("NAME"), "NAME").(string),
			Port: panicIfEmpty(os.Getenv("PORT"), "PORT").(string),
		},
		OrderCloud: OrderCloudConfig{
			ClientID:                 panicIfEmpty(os.Getenv("ORDER_CLOUD_CLIENT_ID"), "ORDER_CLOUD_CLIENT_ID").(string),
			ClientSecret:             panicIfEmpty(os.Getenv("ORDER_CLOUD_CLIENT_SECRET"), "ORDER_CLOUD_CLIENT_SECRET").(string),
			Username:                 panicIfEmpty(os.Getenv("ORDER_CLOUD_USERNAME"), "ORDER_CLOUD_USERNAME").(string),
			Password:                 panicIfEmpty(os.Getenv("ORDER_CLOUD_PASSWORD"), "ORDER_CLOUD_PASSWORD").(string),
			AccessTokenFetchDuration: panicIfEmpty(GetInt(os.Getenv("ORDER_CLOUD_ACCESS_TOKEN_FETCH_DURATION")), "ORDER_CLOUD_ACCESS_TOKEN_FETCH_DURATION").(int),
			OrderCloudEngine:         panicIfEmpty(os.Getenv("ORDER_CLOUD_ENGINE"), "ORDER_CLOUD_ENGINE").(string),
			SellerCenterMiddleware:   panicIfEmpty(os.Getenv("SELLER_CENTER_MIDDLEWARE"), "SELLER_CENTER_MIDDLEWARE").(string),
		},
		DB: DBConfig{
			Host:     panicIfEmpty(os.Getenv("DB_HOST"), "DB_HOST").(string),
			Port:     panicIfEmpty(os.Getenv("DB_PORT"), "DB_PORT").(string),
			Username: panicIfEmpty(os.Getenv("DB_USERNAME"), "DB_USERNAME").(string),
			Password: panicIfEmpty(os.Getenv("DB_PASSWORD"), "DB_PASSWORD").(string),
			Name:     panicIfEmpty(os.Getenv("DB_NAME"), "DB_NAME").(string),
			Schema:   panicIfEmpty(os.Getenv("DB_SCHEMA"), "DB_SCHEMA").(string),
			SSLMode:  panicIfEmpty(os.Getenv("DB_SSL_MODE"), "DB_SSL_MODE").(string),
		},
	}

	return appConfig
}

func GetInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}

func GetBool(val string) bool {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return b
}

func GetNonEmptyData(val interface{}, defaultVal interface{}) interface{} {
	switch val.(type) {
	case string:
		if val.(string) == "" {
			return defaultVal
		}
	case int:
		if val.(int) == 0 {
			return defaultVal
		}
	}

	return val
}

func panicIfEmpty(val interface{}, name string) interface{} {
	switch vv := val.(type) {
	case string:
		if vv == "" {
			panic(name + " config is empty")
		}
	case int:
		if vv == 0 {
			panic(name + " config is empty")
		}
	}

	return val
}
