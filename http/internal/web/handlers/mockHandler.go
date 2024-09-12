package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	mockuser "github.com/Lafetz/loyalty_marketplace/internal/web/mockUser"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const mockID = "4f981b0b-accf-4eb7-8018-7cd651c7e907"

type userSignin struct {
	Username string `json:"username" binding:"required,max=150"`
	Password string `json:"password" binding:"required,min=8,max=500" `
}

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ginUser userSignin
		if err := c.ShouldBindJSON(&ginUser); err != nil {
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
		id, _ := uuid.Parse(mockID)
		user := mockuser.User{
			Id:       id,
			Username: "helloworld",
			Email:    "hellow@world.com",
		}
		token, err := jwtauth.CreateJwt(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "internal server Error",
			})
			return
		}

		c.SetCookie("Authorization", token, 24*60*60, "/", "localhost", true, true)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})

	}
}

func GetTierCustomerHandler(tierSvc *tier.TierSvc, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		md, _ := uuid.Parse("4f981b0b-accf-4eb7-8018-7cd651c7e907")
		cusid, _ := uuid.Parse("4f981b0b-accf-4eb7-8018-7cd651c7e922")

		customerTier, err := tierSvc.GetTierCustomer(c.Request.Context(), md, cusid)
		if err != nil {

			if errors.Is(err, tier.ErrUnauth) {
				c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized access"})
			} else if errors.Is(err, tier.ErrNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"Error": "Customer not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal server error"})
			}
			return
		}

		c.JSON(http.StatusOK, customerTier)
	}
}
