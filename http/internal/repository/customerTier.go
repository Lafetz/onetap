package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/repository/gen"
	"github.com/google/uuid"
)

func (store *Store) CreateCustomerTier(ctx context.Context, customerTier tier.CustomerTier) error {
	err := store.queries.CreateCustomerTier(ctx, gen.CreateCustomerTierParams{
		MerchantID: customerTier.MerchantID,
		CustomerID: customerTier.CustomerID,
		Points:     int32(customerTier.Points),
		TierName:   customerTier.TierName,
	})

	return err
}

func (store *Store) GetCustomerTier(ctx context.Context, merchantId uuid.UUID, customerId uuid.UUID) (tier.CustomerTier, error) {
	ct, err := store.queries.GetCustomerTier(ctx, gen.GetCustomerTierParams{
		MerchantID: merchantId,
		CustomerID: customerId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tier.CustomerTier{}, tier.ErrNotFound
		}
		return tier.CustomerTier{}, err
	}
	return tier.CustomerTier{
		MerchantID: ct.MerchantID,
		CustomerID: ct.CustomerID,
		TierName:   ct.TierName,
		Points:     int(ct.Points),
	}, nil
}

func (store *Store) UpdateCustomerTier(ctx context.Context, updatedCustomerTier tier.CustomerTier) error {
	err := store.queries.UpdateCustomerTier(ctx, gen.UpdateCustomerTierParams{
		MerchantID: updatedCustomerTier.MerchantID,
		CustomerID: updatedCustomerTier.CustomerID,
		TierName:   updatedCustomerTier.TierName,
		Points:     int32(updatedCustomerTier.Points),
	})
	return err
}

func (store *Store) DeleteCustomerTier(ctx context.Context, merchantId uuid.UUID, customerId uuid.UUID) error {
	err := store.queries.DeleteCustomerTier(ctx, gen.DeleteCustomerTierParams{
		MerchantID: merchantId,
		CustomerID: customerId,
	})
	return err
}
