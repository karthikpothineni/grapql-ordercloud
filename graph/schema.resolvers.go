package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"mpmy-product-service/graph/generated"
	"mpmy-product-service/graph/model"
)

// FavoriteProduct is the resolver for the favoriteProduct field.
func (r *mutationResolver) FavoriteProduct(ctx context.Context, productID string, isFavorite bool) (*model.UserProductFavorite, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.FavoriteProduct(userID, productID, isFavorite)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ProductsV2 is the resolver for the productsV2 field.
func (r *queryResolver) ProductsV2(ctx context.Context, catalogID *string, categoryID *string, supplierID *string, isFavorite *bool, search *string, page *string, pageSize *string, sortBy *string, extraFilters map[string]interface{}) (*model.ProductResponseV2, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.GetProductsV2(catalogID, categoryID, supplierID, userID, page, pageSize, sortBy, search, isFavorite, extraFilters, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context, catalogID *string, categoryID *string, supplierID *string, isFavorite *bool, search *string, page *string, pageSize *string, sortBy *string, extraFilters map[string]interface{}) (*model.ProductResponse, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.GetProducts(catalogID, categoryID, supplierID, userID, page, pageSize, sortBy, search, isFavorite, extraFilters, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// SimilarProducts is the resolver for the similarProducts field.
func (r *queryResolver) SimilarProducts(ctx context.Context, productID string, page *string, pageSize *string) (*model.ProductResponse, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.GetSimilarProducts(productID, userID, page, pageSize, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// RecommendProducts is the resolver for the recommendProducts field.
func (r *queryResolver) RecommendProducts(ctx context.Context, productID string, page *string, pageSize *string) (*model.ProductResponseV2, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.GetRecommendProducts(productID, page, pageSize, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.ProductItem, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.GetProduct(id, userID, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ProductV2 is the resolver for the productV2 field.
func (r *queryResolver) ProductV2(ctx context.Context, id string) (*model.LatestProductItems, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.ProductService.GetProductV2(id, userID, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// PriceSchedules is the resolver for the priceSchedules field.
func (r *queryResolver) PriceSchedules(ctx context.Context, productID string, page *string, pageSize *string) (*model.PriceScheduleResponse, error) {
	result, err := r.PriceScheduleService.GetPriceSchedule(productID, page, pageSize, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Categories is the resolver for the categories field.
func (r *queryResolver) Categories(ctx context.Context, catalogID *string, depth *string) (*model.CategoryResponse, error) {
	result, err := r.CategoryService.GetCategories(catalogID, depth, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TrendingProducts is the resolver for the trendingProducts field.
func (r *queryResolver) TrendingProducts(ctx context.Context) (*model.ProductResponse, error) {
	result, err := r.ProductService.GetTrendingProducts(r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetProductFilter is the resolver for the getProductFilter field.
func (r *queryResolver) GetProductFilter(ctx context.Context, search string) ([]*model.ProductFilter, error) {
	result, err := r.ProductService.GetProductFilterMiddleWare(search, r.LoginService.AccessToken)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// RecentSearches is the resolver for the recentSearches field.
func (r *queryResolver) RecentSearches(ctx context.Context, page *string, pageSize *string) ([]*model.RecentSearch, error) {
	userID, err := GetCurrentUserID(ctx)
	if userID == nil {
		return nil, errors.New("invalid user")
	}
	if err != nil {
		return nil, err
	}

	result, err := r.RecentSearchService.GetRecentSearches(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
