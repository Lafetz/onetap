package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateCashback(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store)

	newCashback := cashback.Cashback{
		ID:               uuid.New(),
		MerchantID:       uuid.New(),
		Name:             "Summer Sale",
		Description:      "Discount on all summer products",
		Percentage:       0.5,
		EligibleProducts: []uuid.UUID{},
		Active:           true,
		Expiration:       time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := cashbackSvc.CreateCashback(ctx, newCashback)
	assert.NoError(t, err)
}
func TestGetCashback(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store)

	merchantID := uuid.New()
	cashbackID := uuid.New()

	newCashback := cashback.Cashback{
		ID:               cashbackID,
		MerchantID:       merchantID,
		Name:             "Summer Sale",
		Description:      "Discount on all summer products",
		Percentage:       0.5,
		EligibleProducts: []uuid.UUID{},
		Active:           true,
		Expiration:       time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := cashbackSvc.CreateCashback(ctx, newCashback)
	assert.NoError(t, err)

	cashback, err := cashbackSvc.GetCashback(ctx, cashbackID)
	assert.NoError(t, err)
	assert.Equal(t, newCashback.ID, cashback.ID)
	assert.Equal(t, newCashback.Name, cashback.Name)
}

func TestUpdateCashback(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store)

	merchantID := uuid.New()
	cashbackID := uuid.New()

	newCashback := cashback.Cashback{
		ID:               cashbackID,
		MerchantID:       merchantID,
		Name:             "Summer Sale",
		Description:      "Discount on all summer products",
		Percentage:       0.5,
		EligibleProducts: []uuid.UUID{},
		Active:           true,
		Expiration:       time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := cashbackSvc.CreateCashback(ctx, newCashback)
	assert.NoError(t, err)

	updatedCashback := newCashback
	updatedCashback.Description = "Updated description"
	updatedCashback.UpdatedAt = time.Now()
	err = cashbackSvc.UpdateCashback(ctx, merchantID, updatedCashback)
	assert.NoError(t, err)

	cashback, err := cashbackSvc.GetCashback(ctx, cashbackID)
	assert.NoError(t, err)
	assert.Equal(t, updatedCashback.Description, cashback.Description)
}

func TestDeleteCashback(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store)

	merchantID := uuid.New()
	cashbackID := uuid.New()

	newCashback := cashback.Cashback{
		ID:               cashbackID,
		MerchantID:       merchantID,
		Name:             "Summer Sale",
		Description:      "Discount on all summer products",
		Percentage:       0.5,
		EligibleProducts: []uuid.UUID{},
		Active:           true,
		Expiration:       time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := cashbackSvc.CreateCashback(ctx, newCashback)
	assert.NoError(t, err)

	err = cashbackSvc.DeleteCashback(ctx, merchantID, cashbackID)
	assert.NoError(t, err)

	_, err = cashbackSvc.GetCashback(ctx, cashbackID)
	assert.Error(t, err)
	assert.Equal(t, cashback.ErrNotFound, err)
}

func TestListCashbacks(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store)

	merchantID := uuid.New()

	newCashback := cashback.Cashback{
		ID:               uuid.New(),
		MerchantID:       merchantID,
		Name:             "Winter Sale",
		Description:      "Discount on all winter products",
		Percentage:       0.15,
		EligibleProducts: []uuid.UUID{},
		Active:           true,
		Expiration:       time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	err := cashbackSvc.CreateCashback(ctx, newCashback)
	assert.NoError(t, err)

	cashbacks, err := cashbackSvc.ListCashbacks(ctx, merchantID)
	assert.NoError(t, err)
	assert.Len(t, cashbacks, 1)
	assert.Equal(t, newCashback.ID, cashbacks[0].ID)
}
