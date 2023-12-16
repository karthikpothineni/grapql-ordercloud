package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mpmy-product-service/config"
	"mpmy-product-service/constants"
	"mpmy-product-service/utils"
	"os"
	"strings"
)

func tokenNotPresent(c *gin.Context) {
	c.JSON(401, gin.H{
		"errors": []string{
			"token missing",
		},
		"data": nil,
	})
	c.Abort()
}
func tokenInvalid(c *gin.Context) {
	c.JSON(401, gin.H{
		"errors": []string{
			"invalid token",
		},
		"data": nil,
	})
	c.Abort()
}

func handleToken(c *gin.Context, tokenString string) {
	c.Set(constants.TokenString, tokenString)
	utils.HandlePanic(c)
	ky := os.Getenv("JWT_SECRET")
	var hmacSecret []byte

	if ky != "" {
		hmacSecret = []byte(ky)
	} else {
		panic("SECRET not found")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil || !token.Valid {
		fmt.Printf("token is invalid, error: %v", err)
		tokenInvalid(c)
	} else {
		claims := token.Claims.(jwt.MapClaims)
		for key, val := range claims {
			fmt.Printf("%v: %v\n", key, val)
		}
		c.Set("USER_CLAIM", claims)
		c.Next()
	}
}

func PublicApi(c *gin.Context) bool {
	query := c.Value(constants.RequestBody).(string)
	c1 := strings.Contains(query, "_service")
	return c1
}

func AuthMiddleware(appConfig config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		fullPath := c.FullPath()
		if fullPath == "/health" {
			c.Next()
		} else {
			if (fullPath == "/query") && PublicApi(c) {
				c.Next()
			} else {
				utils.HandlePanic(c)
				tokenInParam := c.Params.ByName("jwt_token")
				utils.Logger(constants.DefaultLogLevel, "tokenInParam")
				utils.Logger(constants.DefaultLogLevel, tokenInParam)
				if tokenInParam != "" {
					tokenNotPresent(c)
				} else {
					tokenStringHeader := c.Request.Header["Authorization"]
					if tokenStringHeader == nil {
						tokenNotPresent(c)
					} else {
						utils.Logger(constants.DefaultLogLevel, "tokenStringHeader")
						utils.Logger(constants.DefaultLogLevel, tokenStringHeader)
						tokenString := tokenStringHeader[0]
						split := strings.Split(tokenString, " ")
						switch len(split) {
						case 1:
							tokenString = split[0]
							break
						case 2:
							tokenString = split[1]
							break
						case 3:
							tokenString = split[2]
							break
						}
						utils.Logger(constants.DefaultLogLevel, "/*tokenString", tokenString, "tokenString*/")
						if tokenString == "" {
							tokenNotPresent(c)
						} else {
							handleToken(c, tokenString)
						}
					}
				}
			}
		}
	}
}
