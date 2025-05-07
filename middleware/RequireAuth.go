package middleware

import (
	"PadelIn/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var CurrentUser models.EndUser

func RequireAuth(c *gin.Context) {
	if os.Getenv("ENABLED_SECURITY") == "NO" {
		c.Next()
		return
	}
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString = tokenString[7:] // Removing "Bearer " from the start of the token

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.EndUser
		user.ID = uint(claims["sub"].(float64))
		user.InstitutionID = uint(claims["inst"].(float64))
		CurrentUser = user
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
