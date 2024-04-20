package middleware

import (
	"fmt"
	"net/http"
	"whispr-golang/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No Authorization header provided")})
			ctx.Abort()
			return 
		}

		claims, err := helpers.ValidateToken(clientToken)

		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("username", claims.Username)
		ctx.Set("uid", claims.Uid)
		ctx.Next()
	}
}