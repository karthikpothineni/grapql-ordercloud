package repository

import (
	"github.com/jmoiron/sqlx"
	"mpmy-product-service/config"
	"mpmy-product-service/graph/model"
	"time"
)

const (
	// RecentSearchesTableName is the name of recent searches table
	RecentSearchesTableName = "recent_searches"
)

type RecentSearchesRepository struct {
	db                 *sqlx.DB
	dbConfig           config.DBConfig
}

func NewRecentSearchesRepository(db *sqlx.DB, dbConfig config.DBConfig) *RecentSearchesRepository {
	return &RecentSearchesRepository{
		db:                 db,
		dbConfig:           dbConfig,
	}
}

func (repo *RecentSearchesRepository) SaveRecentSearch(userId, keyword *string) (*model.RecentSearch, error) {
	currentTime := time.Now().UTC()
	recentSearch := &model.RecentSearch{
		SearchKeyword: keyword,
		UserID: userId,
		CreatedAt: &currentTime,
	}

	query := "INSERT INTO " + repo.dbConfig.Schema + "." + RecentSearchesTableName + "(" +
		"search_keyword, user_id, created_at)" +
		"VALUES(:search_keyword, :user_id, :created_at) " +
		"RETURNING id"

	rows, err := repo.db.NamedQuery(query, recentSearch)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var lastInsertID int
	if rows.Next() {
		err = rows.Scan(&lastInsertID)
		if err != nil {
			return nil, err
		}
	}

	recentSearch.ID = &lastInsertID

	return recentSearch, nil
}

func (repo *RecentSearchesRepository) GetRecentSearches(userID, page, pageSize *string) ([]*model.RecentSearch, error) {
	var recentSearches []*model.RecentSearch
	var err error

	query := "SELECT * FROM " + repo.dbConfig.Schema + "." + RecentSearchesTableName + " WHERE user_id = $1 ORDER BY created_at DESC"

	if page != nil && pageSize != nil {
		query += " LIMIT $2 OFFSET $3"

		limit := GetIntFromStringPointer(pageSize)
		offset := (GetIntFromStringPointer(page) * limit) - limit

		err = repo.db.Select(&recentSearches, query, userID, limit, offset)
	} else {
		err = repo.db.Select(&recentSearches, query, userID)
	}

	if err != nil {
		return nil, err
	}

	return recentSearches, nil
}