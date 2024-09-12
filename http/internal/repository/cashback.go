package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/repository/gen"
	"github.com/google/uuid"
)

func (store *Store) CreateCashback(ctx context.Context, cashback cashback.Cashback) error {
	err := store.queries.CreateCashback(ctx, gen.CreateCashbackParams{
		ID:               cashback.ID,
		MerchantID:       cashback.MerchantID,
		Name:             cashback.Name,
		Description:      cashback.Description,
		Percentage:       cashback.Percentage,
		EligibleProducts: cashback.EligibleProducts,
		Active:           cashback.Active,
		Expiration:       cashback.Expiration,
	})
	return err
}

func (store *Store) GetCashback(ctx context.Context, cashbackID uuid.UUID) (cashback.Cashback, error) {
	cb, err := store.queries.GetCashback(ctx, cashbackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return cashback.Cashback{}, cashback.ErrNotFound
		}
		return cashback.Cashback{}, err
	}
	return cashback.Cashback{
		ID:               cb.ID,
		MerchantID:       cb.MerchantID,
		Name:             cb.Name,
		Description:      cb.Description,
		Percentage:       cb.Percentage,
		EligibleProducts: cb.EligibleProducts,
		Active:           cb.Active,
		Expiration:       cb.Expiration,
		CreatedAt:        cb.CreatedAt,
		UpdatedAt:        cb.UpdatedAt,
	}, nil
}

func (store *Store) UpdateCashback(ctx context.Context, cashback cashback.Cashback) error {
	err := store.queries.UpdateCashback(ctx, gen.UpdateCashbackParams{
		ID:               cashback.ID,
		Name:             cashback.Name,
		Description:      cashback.Description,
		Percentage:       cashback.Percentage,
		EligibleProducts: cashback.EligibleProducts,
		Active:           cashback.Active,
		Expiration:       cashback.Expiration,
	})
	return err
}

func (store *Store) DeleteCashback(ctx context.Context, cashbackID uuid.UUID) error {
	err := store.queries.DeleteCashback(ctx, cashbackID)
	return err
}

func (store *Store) ListCashbacks(ctx context.Context, merchantID uuid.UUID) ([]cashback.Cashback, error) {
	cbs, err := store.queries.ListCashbacks(ctx, merchantID)
	if err != nil {
		return nil, err
	}
	if len(cbs) == 0 {
		return []cashback.Cashback{}, nil
	}
	var cashbacks []cashback.Cashback
	for _, cb := range cbs {
		cashbacks = append(cashbacks, cashback.Cashback{
			ID:               cb.ID,
			MerchantID:       cb.MerchantID,
			Name:             cb.Name,
			Description:      cb.Description,
			Percentage:       cb.Percentage,
			EligibleProducts: cb.EligibleProducts,
			Active:           cb.Active,
			Expiration:       cb.Expiration,
			CreatedAt:        cb.CreatedAt,
			UpdatedAt:        cb.UpdatedAt,
		})
	}
	return cashbacks, nil
}
