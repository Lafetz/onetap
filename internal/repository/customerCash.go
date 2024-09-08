package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/repository/gen"
	"github.com/google/uuid"
)

func (store *Store) CreateCashbackUser(ctx context.Context, cashbackUser cashback.CashbackUser) error {
	err := store.queries.CreateCashbackUser(ctx, gen.CreateCashbackUserParams{
		MerchantID: cashbackUser.MerchantID,
		CashbackID: cashbackUser.CashbackID,
		UserID:     cashbackUser.UserID,
		Points:     float64(cashbackUser.Points),
	})
	return err
}

func (store *Store) GetCashbackUser(ctx context.Context, userID uuid.UUID, cashbackID uuid.UUID) (cashback.CashbackUser, error) {
	cu, err := store.queries.GetCashbackUser(ctx, gen.GetCashbackUserParams{
		UserID:     userID,
		CashbackID: cashbackID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cashback.CashbackUser{}, cashback.ErrNotFound
		}
		return cashback.CashbackUser{}, err
	}
	return cashback.CashbackUser{
		CashbackID: cu.CashbackID,
		UserID:     cu.UserID,
		Points:     float64(cu.Points),
	}, nil
}

func (store *Store) UpdateCashbackUser(ctx context.Context, cashbackUser cashback.CashbackUser) error {
	err := store.queries.UpdateCashbackUser(ctx, gen.UpdateCashbackUserParams{
		CashbackID: cashbackUser.CashbackID,
		UserID:     cashbackUser.UserID,
		Points:     float64(cashbackUser.Points),
	})
	return err
}

func (store *Store) DeleteCashbackUser(ctx context.Context, userID uuid.UUID, cashbackID uuid.UUID) error {
	err := store.queries.DeleteCashbackUser(ctx, gen.DeleteCashbackUserParams{
		UserID:     userID,
		CashbackID: cashbackID,
	})
	return err
}

func (store *Store) ListCashbackUsers(ctx context.Context, merchantID uuid.UUID) ([]cashback.CashbackUser, error) {
	cus, err := store.queries.ListCashbackUsers(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	var cashbackUsers []cashback.CashbackUser
	for _, cu := range cus {
		cashbackUsers = append(cashbackUsers, cashback.CashbackUser{
			CashbackID: cu.CashbackID,
			UserID:     cu.UserID,
			Points:     float64(cu.Points),
		})
	}
	return cashbackUsers, nil
}
