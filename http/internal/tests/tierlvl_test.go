package integration_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTierLevel(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())

	merchantId := uuid.New()

	err := tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Platinum",
		MinPoints:  2000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	assert.NoError(t, err)

	tier, err := tierSvc.GetTierLevel(ctx, merchantId, "Platinum")
	assert.NoError(t, err)
	assert.Equal(t, "Platinum", tier.Name)
	assert.Equal(t, 2000, tier.MinPoints)
}
func TestUpdateTierLevel(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())

	merchantId := uuid.New()
	tierlvl := tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Gold",
		MinPoints:  1500,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := tierSvc.CreateTierLevel(ctx, tierlvl)
	//

	//
	assert.NoError(t, err)
	tierlvl.Name = "Diamond"
	tierlvl.MinPoints = 3000
	err = tierSvc.UpdateTierLevel(ctx, tierlvl)
	assert.NoError(t, err)
	tier, err := tierSvc.GetTierLevel(ctx, merchantId, "Diamond")
	assert.NoError(t, err)
	assert.Equal(t, "Diamond", tier.Name)
	assert.Equal(t, 3000, tier.MinPoints)
}
func TestRemoveTierLevel(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store, &mockNot{}, slog.Default())

	merchantId := uuid.New()

	err := tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Silver",
		MinPoints:  1000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	assert.NoError(t, err)

	err = tierSvc.DeleteTierLevel(ctx, merchantId, "Silver")
	assert.NoError(t, err)

	_, err = tierSvc.GetTierLevel(ctx, merchantId, "Silver")
	assert.Error(t, err)
	assert.Equal(t, tier.ErrNotFound, err)
}
