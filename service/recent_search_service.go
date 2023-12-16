package service

import (
	"github.com/jmoiron/sqlx"
	"mpmy-product-service/config"
	"mpmy-product-service/graph/model"
	"mpmy-product-service/repository"
)

type RecentSearchService struct {
	recentSearchesRepo   *repository.RecentSearchesRepository
}

func NewRecentSearchService(db *sqlx.DB, dbConfig config.DBConfig) *RecentSearchService {
	return &RecentSearchService{
		recentSearchesRepo: repository.NewRecentSearchesRepository(db, dbConfig),
	}
}

func (s *RecentSearchService) GetRecentSearches(userID, page, pageSize *string) ([]*model.RecentSearch, error) {
	recentSearches, err := s.recentSearchesRepo.GetRecentSearches(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return recentSearches, nil
}