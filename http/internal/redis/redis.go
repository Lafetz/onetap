package pubsub

import (
	"context"
	"encoding/json"
	"log/slog"

	core "github.com/Lafetz/loyalty_marketplace/internal/loyalty"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

const (
	PUB_TOPIC = "notification"
	SUB_TOPIC = "order"
)

type Order struct {
	CustomerID uuid.UUID `json:"customerId"`
	MerchantID uuid.UUID `json:"merchantId"`
	ProductID  uuid.UUID `json:"productId"`
	Amount     int       `json:"amount"`
}
type Pubsub struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func (p *Pubsub) ReceiveMessage(tierSvc *tier.TierSvc, cashbackSvc *cashback.CashbackSvc) {
	subscriber := p.redisClient.Subscribe(context.TODO(), SUB_TOPIC)
	p.logger.Debug("started to receive msgs from redis")
	for {
		order := Order{}

		msg, err := subscriber.ReceiveMessage(context.TODO())
		p.logger.Debug("received message from reds", "msg", msg)
		if err != nil {
			p.logger.Error("error reading msg from redis", "error", err.Error())
		}
		if err := json.Unmarshal([]byte(msg.Payload), &order); err != nil {
			p.logger.Error("error unmarshalling payload", "error", err.Error())
		}
		err = tierSvc.ProcessOrder(context.TODO(), order.CustomerID, order.MerchantID, order.Amount)
		if err != nil {
			p.logger.Error("error processing order tierSvc", "error", err.Error())
		}
		p.logger.Debug("finished processing order for tier svc")
		err = cashbackSvc.ProcessOrder(context.TODO(), order.CustomerID, order.MerchantID, order.ProductID, float64(order.Amount))

		if err != nil {
			p.logger.Error("error processing order cashbackSvc", "error", err.Error())
		}
		p.logger.Debug("finished processing order for cashback svc")
		p.logger.Debug("finished processing message")
	}
}
func (p *Pubsub) SendNotification(ctx context.Context, noti core.Notification) {
	p.logger.Debug("received request to send notification", "notification", noti)
	payload, err := json.Marshal(noti)
	if err != nil {
		p.logger.Error("error marshalling payload", "error", err.Error())
	}
	if err := p.redisClient.Publish(ctx, PUB_TOPIC, payload).Err(); err != nil {
		p.logger.Error("error publishing payload to redis", "error", err.Error())
		return
	}
	p.logger.Debug("published notification", "notification", noti)
}
func NewPubsub(addr, password string, logger *slog.Logger) *Pubsub {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &Pubsub{
		redisClient: redisClient,
		logger:      logger,
	}
}
