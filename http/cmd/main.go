package main

import (
	"log"

	"github.com/Lafetz/loyalty_marketplace/internal/config"
	customlogger "github.com/Lafetz/loyalty_marketplace/internal/logger"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	pubsub "github.com/Lafetz/loyalty_marketplace/internal/redis"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/Lafetz/loyalty_marketplace/internal/web"
)

const (
	serviceName = "loyality"
	version     = "1.0.0"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}
	logger := customlogger.NewLogger(cfg.Env, cfg.LogLevel, serviceName, version)
	if err != nil {
		log.Fatal(err)
	}
	pubsub := pubsub.NewPubsub(cfg.RedisUrl, cfg.RedisPass, logger)

	db, err := repository.OpenDB(cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	store := repository.NewDb(db)
	tierSvc := tier.NewTierSvc(store, store, pubsub, logger)
	cashbackSvc := cashback.NewCashbackSvc(store, store, pubsub, logger)

	srv := web.NewApp(tierSvc, cashbackSvc, 8080, logger)

	go pubsub.ReceiveMessage(tierSvc, cashbackSvc)
	srv.Run()

}
