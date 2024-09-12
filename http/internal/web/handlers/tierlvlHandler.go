package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TierReq struct {
	Name     string `json:"name" binding:"required,min=3"`
	MinPoint int    `json:"minpoint" binding:"required"`
}

func CreateTier(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		id, err := uuid.Parse(user.(*jwtauth.UserToken).Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "internal server Error",
			})
			return
		}
		var tierReq TierReq
		if err := c.ShouldBind(&tierReq); err != nil {
			_, ok := err.(validator.ValidationErrors)

			if ok {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"Errors": ValidateModel(err),
				})
				return

			}
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Error processing request body",
			})
			return
		}
		t := tier.NewTierLevel(id, tierReq.Name, tierReq.MinPoint)
		err = tierSvc.CreateTierLevel(c, t)
		if err != nil {
			if errors.Is(err, tier.ErrDepulicateTier) {
				c.JSON(http.StatusConflict, gin.H{"Error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "internal server Error",
			})
			logger.Error(err.Error())
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "tier added",
			"tier":    t,
		})

	}
}
func GetTier(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		merchantIdStr := c.Param("merchantId")
		name := c.Param("name")

		merchantId, err := uuid.Parse(merchantIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid merchant ID",
			})
			return
		}

		tier, err := tierSvc.GetTierLevel(c, merchantId, name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"Error": "Tier not found",
			})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tier": tier,
		})
	}
}

func DeleteTier(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		id, err := uuid.Parse(user.(*jwtauth.UserToken).Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "internal server Error",
			})
			return
		}
		merchantIdStr := c.Param("merchantId")
		name := c.Param("name")

		merchantId, err := uuid.Parse(merchantIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid merchant ID",
			})
			return
		}
		//
		if id != merchantId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		_, err = tierSvc.GetTierLevel(c, merchantId, name)
		if err != nil {

			c.JSON(http.StatusNotFound, gin.H{
				"Error": "Tier not found",
			})
			return
		}
		err = tierSvc.DeleteTierLevel(c, merchantId, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal server error",
			})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Tier deleted",
		})
	}
}

func ListTier(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantIdStr := c.Param("merchantId")

		merchantId, err := uuid.Parse(merchantIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid merchant ID",
			})
			return
		}

		tiers, err := tierSvc.ListTierLevels(c, merchantId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal server error",
			})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tiers": tiers,
		})
	}
}

type TierUpdateReq struct {
	Name     string `json:"name" binding:"required,min=3"`
	MinPoint int    `json:"minpoint" binding:"required"`
}

func UpdateTierLevelHandler(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {

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
		var req TierUpdateReq
		if err := c.ShouldBind(&req); err != nil {
			_, ok := err.(validator.ValidationErrors)

			if ok {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"Errors": ValidateModel(err),
				})
				return

			}
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "Error processing request body",
			})
			logger.Error(err.Error())
			return
		}
		name := c.Param("name")
		t, err := tierSvc.GetTierLevel(c, merchantId, name)
		if err != nil {

			c.JSON(http.StatusNotFound, gin.H{
				"Error": "Tier not found",
			})
			return
		}
		if t.MerchantID != merchantId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}
		println(req.Name, req.MinPoint)
		err = tierSvc.UpdateTierLevel(c, tier.TierLevel{
			ID:         t.ID,
			MerchantID: merchantId,
			Name:       req.Name,
			MinPoints:  req.MinPoint,
		})
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tier level updated successfully"})
	}
}
