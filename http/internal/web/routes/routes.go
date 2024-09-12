package routes

import (
	"log/slog"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/web/handlers"
	"github.com/Lafetz/loyalty_marketplace/internal/web/middleware"
	"github.com/gin-gonic/gin"
)

func InitAppRoutes(gin *gin.Engine, tierSvc *tier.TierSvc, cashbackSvc *cashback.CashbackSvc, logger *slog.Logger) {
	// url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")

	gin.POST("/v1/merchants/:merchantId/tiers", middleware.RequireAuth(), handlers.CreateTier(tierSvc, logger))
	gin.GET("/v1/merchants/:merchantId/tiers", handlers.ListTier(tierSvc, logger))
	gin.GET("/v1/merchants/:merchantId/tiers/:name", handlers.GetTier(tierSvc, logger)) //
	gin.PUT("/v1/merchants/:merchantId/tiers/:name", middleware.RequireAuth(), handlers.UpdateTierLevelHandler(tierSvc, logger))
	gin.DELETE("/v1/merchants/:merchantId/tiers/:name", middleware.RequireAuth(), handlers.DeleteTier(tierSvc, logger))
	gin.POST("/v1/merchants/customers/:customerId/promote", handlers.PromoteCustomerTier(tierSvc, logger))
	gin.POST("/v1/merchants/customers/:customerId/demote", handlers.DemoteCustomerTier(tierSvc, logger))
	// Cashback Routes
	gin.POST("/v1/cashbacks", middleware.RequireAuth(), handlers.CreateCashback(cashbackSvc, logger))
	gin.PUT("/v1/cashbacks/:cashbackID", middleware.RequireAuth(), handlers.UpdateCashback(cashbackSvc, logger))
	gin.DELETE("/v1/cashbacks/:cashbackID", middleware.RequireAuth(), handlers.DeleteCashback(cashbackSvc, logger))
	gin.GET("/v1/cashbacks/:cashbackID", handlers.ListCashbacks(cashbackSvc, logger))
	gin.GET("/v1/cashbacks", handlers.ListCashbacks(cashbackSvc, logger))
	//
	//for testing only
	gin.GET("/v1/merchants/:merchantId/:customerId", handlers.GetTierCustomerHandler(tierSvc, logger))
	gin.POST("/v1/signin", handlers.Signin())

}
