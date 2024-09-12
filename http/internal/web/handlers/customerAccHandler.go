package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PromoteCustomerTier(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		merchantId, err := uuid.Parse(user.(*jwtauth.UserToken).Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid merchant ID",
			})
			return
		}
		customerId, err := uuid.Parse(c.Param("customerId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid customer ID"})
			return
		}

		err = tierSvc.PromoteCustomerTier(c.Request.Context(), merchantId, customerId)
		if err != nil {
			if errors.Is(err, tier.ErrNoTier) {
				c.JSON(http.StatusNotFound, gin.H{"Error": "Tier not found"})
			} else if errors.Is(err, tier.ErrHighTier) {
				c.JSON(http.StatusConflict, gin.H{"Error": "Customer is already at the highest tier"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer tier promoted successfully"})
	}
}
func DemoteCustomerTier(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		merchantId, err := uuid.Parse(user.(*jwtauth.UserToken).Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid merchant ID",
			})
			return
		}

		customerId, err := uuid.Parse(c.Param("customerId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid customer ID"})
			return
		}

		err = tierSvc.DemoteCustomerTier(c.Request.Context(), merchantId, customerId)
		if err != nil {
			if errors.Is(err, tier.ErrNoTier) {
				c.JSON(http.StatusNotFound, gin.H{"Error": "customer is not on any tier"})
			} else if errors.Is(err, tier.ErrLowTier) {
				c.JSON(http.StatusConflict, gin.H{"Error": "Customer is already at the lowest tier"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer tier demoted successfully"})
	}
}
