package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type Order struct {
	CustomerID uuid.UUID `json:"customerId"`
	MerchantID uuid.UUID `json:"merchantId"`
	Amount     int       `json:"amount"`
}

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
})

func main() {
	db, err := repository.OpenDB("postgresql://user:password@localhost:5432/loyality?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	store := repository.NewDb(db)
	svc := tier.NewTierSvc(store, store)

	subscriber := redisClient.Subscribe("send-user-data")
	order := Order{}

	for {
		msg, err := subscriber.ReceiveMessage()
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal([]byte(msg.Payload), &order); err != nil {
			panic(err)
		}
		fmt.Println("Received message from " + msg.Channel + " channel.")
		fmt.Printf("%+v\n", order)
		//
		svc.ProcessOrder(context.TODO(), order.CustomerID, order.MerchantID, order.Amount)

		// ...
	}
}
