package cashback

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (c *CashbackSvc) ProcessOrder(ctx context.Context, userID uuid.UUID, merchantID uuid.UUID, productID uuid.UUID, amount float64) error {

	cashbacks, err := c.cashbackRepo.ListCashbacks(ctx, merchantID)
	if err != nil {
		return err
	}

	for _, cashback := range cashbacks {
		if !cashback.Active || time.Now().After(cashback.Expiration) {
			continue
		}
		eligibleProducts := cashback.EligibleProducts

		if !contains(eligibleProducts, productID) {

			break
		}

		points := cashback.Percentage * amount

		userAcc, err := c.cashbackUserRepo.GetCashbackUser(ctx, userID, cashback.ID)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				userAcc = CashbackUser{
					UserID:     userID,
					CashbackID: cashback.ID,
					Points:     points,
				}
				err := c.cashbackUserRepo.CreateCashbackUser(ctx, userAcc)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			userAcc.Points += points
			err := c.cashbackUserRepo.UpdateCashbackUser(ctx, userAcc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func contains(slice []uuid.UUID, item uuid.UUID) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
func (c *CashbackSvc) CreateCashbackUser(ctx context.Context, cashbackUser CashbackUser) error {
	return c.cashbackUserRepo.CreateCashbackUser(ctx, cashbackUser)
}

func (c *CashbackSvc) GetCashbackUser(ctx context.Context, userID uuid.UUID, cashbackID uuid.UUID) (CashbackUser, error) {
	return c.cashbackUserRepo.GetCashbackUser(ctx, userID, cashbackID)
}

func (c *CashbackSvc) UpdateCashbackUser(ctx context.Context, cashbackUser CashbackUser) error {
	return c.cashbackUserRepo.UpdateCashbackUser(ctx, cashbackUser)
}

func (c *CashbackSvc) DeleteCashbackUser(ctx context.Context, userID uuid.UUID, cashbackID uuid.UUID) error {
	return c.cashbackUserRepo.DeleteCashbackUser(ctx, userID, cashbackID)
}

func (c *CashbackSvc) ListCashbackUsers(ctx context.Context, merchantID uuid.UUID) ([]CashbackUser, error) {
	return c.cashbackUserRepo.ListCashbackUsers(ctx, merchantID)
}
