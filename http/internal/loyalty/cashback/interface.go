package cashback

import (
	"context"

	core "github.com/Lafetz/loyalty_marketplace/internal/loyalty"
	"github.com/google/uuid"
)

type CashBackRepository interface {
	CreateCashback(ctx context.Context, cashback Cashback) error
	UpdateCashback(ctx context.Context, cashback Cashback) error
	GetCashback(ctx context.Context, cashbackID uuid.UUID) (Cashback, error)
	DeleteCashback(ctx context.Context, cashbackID uuid.UUID) error
	ListCashbacks(ctx context.Context, merchantID uuid.UUID) ([]Cashback, error)
}

type CashbackCustRepository interface {
	CreateCashbackUser(ctx context.Context, user CashbackUser) error
	UpdateCashbackUser(ctx context.Context, user CashbackUser) error
	GetCashbackUser(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID) (CashbackUser, error)
	DeleteCashbackUser(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID) error
	ListCashbackUsers(ctx context.Context, merchantID uuid.UUID) ([]CashbackUser, error)
}
type notification interface {
	SendNotification(ctx context.Context, noti core.Notification)
}
