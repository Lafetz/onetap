package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/repository/gen"
	"github.com/google/uuid"
)

func (store *Store) CreateTierLevel(ctx context.Context, tierLevel tier.TierLevel) error {
	err := store.queries.CreateTierLevel(ctx, gen.CreateTierLevelParams{
		TierID:     tierLevel.ID,
		MerchantID: tierLevel.MerchantID,
		Name:       tierLevel.Name,
		MinPoints:  int32(tierLevel.MinPoints),
	})
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "tier_level_merchant_id_name_key"`:
			return tier.ErrDepulicateTier
		default:
			return err
		}
	}
	return err
}

// "tier_level_merchant_id_name_key"
func (store *Store) GetTierLevel(ctx context.Context, merchantid uuid.UUID, name string) (tier.TierLevel, error) {
	t, err := store.queries.GetTierLevel(ctx, gen.GetTierLevelParams{MerchantID: merchantid, Name: name})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tier.TierLevel{}, tier.ErrNotFound
		}
		return tier.TierLevel{}, err
	}
	return tier.TierLevel{
		MerchantID: t.MerchantID,
		Name:       t.Name,
		MinPoints:  int(t.MinPoints),
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.CreatedAt,
	}, nil
}
func (store *Store) UpdateTierLevel(ctx context.Context, updatedTierLevel tier.TierLevel) error {
	err := store.queries.UpdateTierLevel(ctx, gen.UpdateTierLevelParams{
		TierID:    updatedTierLevel.ID,
		Name:      updatedTierLevel.Name,
		MinPoints: int32(updatedTierLevel.MinPoints),
	})
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "tier_level_merchant_id_name_key"`:
			return tier.ErrDepulicateTier
		default:
			return err
		}
	}
	return err
}

func (store *Store) DeleteTierLevel(ctx context.Context, merchantId uuid.UUID, name string) error {
	err := store.queries.DeleteTierLevel(ctx, gen.DeleteTierLevelParams{
		MerchantID: merchantId,
		Name:       name,
	})
	return err
}

func (store *Store) ListTierLevel(ctx context.Context, merchantId uuid.UUID) ([]tier.TierLevel, error) {
	tierLevels, err := store.queries.ListTierLevels(ctx, merchantId)
	if err != nil {
		return nil, err
	}
	if len(tierLevels) == 0 {
		return []tier.TierLevel{}, nil
	}
	var result []tier.TierLevel
	for _, t := range tierLevels {

		result = append(result, tier.TierLevel{
			MerchantID: t.MerchantID,
			Name:       t.Name,
			MinPoints:  int(t.MinPoints),
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		})
	}

	return result, nil
}
