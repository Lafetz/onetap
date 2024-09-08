package cashback

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (c *CashbackSvc) CreateCashback(ctx context.Context, newCashback Cashback) error {

	return c.cashbackRepo.CreateCashback(ctx, newCashback)
}
func (c *CashbackSvc) GetCashback(ctx context.Context, cashbackID uuid.UUID) (Cashback, error) {
	return c.cashbackRepo.GetCashback(ctx, cashbackID)
}
func (c *CashbackSvc) UpdateCashback(ctx context.Context, merchantID uuid.UUID, cashback Cashback) error {
	existingCashback, err := c.cashbackRepo.GetCashback(ctx, cashback.ID)
	if err != nil {
		return err
	}
	if existingCashback.MerchantID != merchantID {
		return ErrNotAuth
	}
	cashback.UpdatedAt = time.Now()

	return c.cashbackRepo.UpdateCashback(ctx, cashback)
}

func (c *CashbackSvc) DeleteCashback(ctx context.Context, merchantID uuid.UUID, cashbackID uuid.UUID) error {
	existingCashback, err := c.cashbackRepo.GetCashback(ctx, cashbackID)
	if err != nil {
		return err
	}
	if existingCashback.MerchantID != merchantID {
		return ErrNotAuth
	}
	return c.cashbackRepo.DeleteCashback(ctx, cashbackID)
}
func (c *CashbackSvc) ListCashbacks(ctx context.Context, merchantID uuid.UUID) ([]Cashback, error) {
	return c.cashbackRepo.ListCashbacks(ctx, merchantID)
}
