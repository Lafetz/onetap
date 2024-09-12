// @/pub/main.go
package main

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const SUB_TOPIC = "order"

type Order struct {
	CustomerID uuid.UUID `json:"customerId"`
	MerchantID uuid.UUID `json:"merchantId"`
	ProductID  uuid.UUID `json:"productId"`
	Amount     int       `json:"amount"`
}

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "cache:6379",
	Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
})

func main() {
	app := fiber.New()
	md, _ := uuid.Parse("4f981b0b-accf-4eb7-8018-7cd651c7e907")
	cusid, _ := uuid.Parse("4f981b0b-accf-4eb7-8018-7cd651c7e922")
	app.Post("/order", func(c *fiber.Ctx) error {
		user := Order{
			CustomerID: cusid,
			ProductID:  cusid,
			MerchantID: md,
			Amount:     50,
		}

		payload, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		if err := redisClient.Publish(ctx, SUB_TOPIC, payload).Err(); err != nil {
			panic(err)
		}

		return c.SendStatus(200)
	})

	app.Listen(":3000")
}
