package handlers

import (
	"net/http"

	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func extractMerchantID(c *gin.Context) (uuid.UUID, bool) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized",
		})
		return uuid.Nil, false
	}

	merchantId, err := uuid.Parse(user.(*jwtauth.UserToken).Id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "Unauthorized",
		})
		return uuid.Nil, false
	}

	return merchantId, true
}
