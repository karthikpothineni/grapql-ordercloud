package db

import (
	"fmt"
	"mpmy-product-service/config"
	"time"

	// postgres db driver
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const DriverName = "postgres"

func Init(dbConfig config.DBConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

	client, err := sqlx.Open(DriverName, dataSource)
	if err != nil {
		return nil, err
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(100)
	client.SetMaxIdleConns(100)

	return client, nil
}
