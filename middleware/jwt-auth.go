package middleware

import (
	"github.com/ByteAKA8bit/golang_api/helper"
	"github.com/ByteAKA8bit/golang_api/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthorizeJWT(service service.JWTService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("failed to process request", "No token found", nil)
			context.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
		token, err := service.ValidateToken(authHeader)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim[user_id]: ", claims["user_id"])
			log.Println("Claim[issuer]: ", claims["issuer"])
		} else {
			response := helper.BuildErrorResponse("token is invalid", err.Error(), nil)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
	}
}
