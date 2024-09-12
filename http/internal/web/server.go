package web

import (
	"log/slog"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/web/routes"
	"github.com/gin-gonic/gin"
)

type App struct {
	gin         *gin.Engine
	port        int
	logger      *slog.Logger
	tierSvc     *tier.TierSvc
	cashbackSvc *cashback.CashbackSvc
}

func NewApp(tierSvc *tier.TierSvc, cashbackSvc *cashback.CashbackSvc, port int, logger *slog.Logger) *App {
	a := &App{
		gin:         gin.Default(),
		port:        port,
		tierSvc:     tierSvc,
		logger:      logger,
		cashbackSvc: cashbackSvc,
	}
	a.gin.Use(corsMiddleware())
	routes.InitAppRoutes(a.gin, a.tierSvc, cashbackSvc, a.logger)
	return a
}
func (a *App) Run() error {
	return a.gin.Run()
}
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
