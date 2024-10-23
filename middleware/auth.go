package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		errEnv := godotenv.Load()

		if errEnv != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": http.StatusText(http.StatusUnauthorized)})
			return
		}

		secret := os.Getenv("SECRET")
		author := ctx.Request.Header.Get("Authorization")
		tokenString := strings.Replace(author, "Bearer ", "", -1)

		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": http.StatusText(http.StatusUnauthorized)})
			return
		}

		token, errTok := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != t.Method {
				return nil, errors.New("methode not metch")
			}
			return []byte(secret), nil
		})

		if errTok != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": http.StatusText(http.StatusForbidden)})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !token.Valid || !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": http.StatusText(http.StatusForbidden)})
			return
		}

		exp, _ := claims["exp"].(float64)

		if time.Now().After(time.Unix(int64(exp), 0)) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Status": "Token Expired"})
			return
		}

		role, _ := claims["Role"].(string)

		if role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": http.StatusText(http.StatusUnauthorized)})
			return
		}

		id, _ := claims["Uid"].(int)
		ctx.Set("id", id)
		ctx.Next()
	}
}
