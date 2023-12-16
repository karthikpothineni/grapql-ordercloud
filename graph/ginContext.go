package graph

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetCurrentUserID(ctx context.Context) (*string, error) {
	gc, err := GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	claims := gc.Value("USER_CLAIM").(jwt.MapClaims)

	id := strconv.FormatFloat(claims["id"].(float64), 'f', 0, 64)
	return &id, nil
}
