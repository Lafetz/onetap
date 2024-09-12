package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CashbackReq struct {
	Name             string      `json:"name" binding:"required"`
	Description      string      `json:"description" binding:"required"`
	Percentage       float64     `json:"percentage" binding:"required,gt=0,lt=1"`
	Expiration       time.Time   `json:"expiration" binding:"required"`
	EligibleProducts []uuid.UUID `json:"eligibleProducts" binding:"required"`
	Active           bool        `json:"active" binding:"required"`
}

func CreateCashback(cashbackSvc *cashback.CashbackSvc, logger *slog.Logger) gin.HandlerFunc {
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
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			return
		}

		var newCashback CashbackReq
		if err := c.BindJSON(&newCashback); err != nil {
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

		if newCashback.Expiration.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Expiration date is required"})
			return
		}

		err = cashbackSvc.CreateCashback(c.Request.Context(), cashback.Cashback{
			ID:               uuid.New(),
			MerchantID:       merchantId,
			Name:             newCashback.Name,
			Description:      newCashback.Description,
			Percentage:       newCashback.Percentage,
			EligibleProducts: newCashback.EligibleProducts,
			Active:           true,
			Expiration:       newCashback.Expiration,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to create cashback"})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cashback created successfully"})
	}
}

func GetCashback(cashbackSvc *cashback.CashbackSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		cashbackIDStr := c.Param("cashbackID")
		cashbackID, err := uuid.Parse(cashbackIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid cashback ID"})
			return
		}

		cashback, err := cashbackSvc.GetCashback(c.Request.Context(), cashbackID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Cashback not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"cashback": cashback})
	}
}
func UpdateCashback(cashbackSvc *cashback.CashbackSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantID, valid := extractMerchantID(c)
		if !valid {
			return
		}
		cashbackIDStr := c.Param("cashbackID")
		cashbackID, err := uuid.Parse(cashbackIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid cashback ID"})
			return
		}
		cashback, err := cashbackSvc.GetCashback(c.Request.Context(), cashbackID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"Error": "Cashback not found"})
			return
		}
		var cashbackreq CashbackReq
		if err := c.BindJSON(&cashbackreq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request payload"})
			return
		}
		cashback.Name = cashbackreq.Name
		cashback.Description = cashbackreq.Description
		cashback.Percentage = cashbackreq.Percentage
		cashback.EligibleProducts = cashbackreq.EligibleProducts
		cashback.Active = cashbackreq.Active
		err = cashbackSvc.UpdateCashback(c.Request.Context(), merchantID, cashback)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to update cashback"})
			logger.Error(err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cashback updated successfully"})
	}
}
func DeleteCashback(cashbackSvc *cashback.CashbackSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		merchantID, valid := extractMerchantID(c)
		if !valid {

			return
		}

		cashbackIDStr := c.Param("cashbackID")
		cashbackID, err := uuid.Parse(cashbackIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid cashback ID"})
			return
		}

		err = cashbackSvc.DeleteCashback(c.Request.Context(), merchantID, cashbackID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete cashback"})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Cashback deleted successfully"})
	}
}

func ListCashbacks(cashbackSvc *cashback.CashbackSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantIDStr := c.Query("merchantID")
		if merchantIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "merchantID is required"})
			return
		}
		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid merchant ID"})
			return
		}
		cashbacks, err := cashbackSvc.ListCashbacks(c.Request.Context(), merchantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve cashbacks"})
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"cashbacks": cashbacks})
	}
}
