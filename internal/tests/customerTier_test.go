package integration_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/tier"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testDbInstance *sql.DB

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestProcessOrder_ExistingCustomer(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store)

	merchantId := uuid.New()
	customerId := uuid.New()

	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Bronze",
		MinPoints:  50,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Gold",
		MinPoints:  1500,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Silver",
		MinPoints:  1000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "Green",
		MinPoints:  500,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	err := tierSvc.ProcessOrder(ctx, customerId, merchantId, 50)
	assert.NoError(t, err)
	customerAcc, err := tierSvc.GetTierCustomer(ctx, merchantId, customerId)
	assert.NoError(t, err)
	assert.Equal(t, 50, customerAcc.Points)
	assert.Equal(t, "Bronze", customerAcc.TierName)

	err = tierSvc.ProcessOrder(ctx, customerId, merchantId, 1050)
	assert.NoError(t, err)

	customerAcc, err = tierSvc.GetTierCustomer(ctx, merchantId, customerId)
	assert.NoError(t, err)
	assert.Equal(t, 1100, customerAcc.Points)
	assert.Equal(t, "Silver", customerAcc.TierName)
}

func TestProcessOrder_MerchantWithoutTiers(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store)

	merchantId := uuid.New()
	customerId := uuid.New()
	err := tierSvc.ProcessOrder(ctx, customerId, merchantId, 150)

	assert.NoError(t, err)

	_, err = tierSvc.GetTierCustomer(ctx, merchantId, customerId)
	assert.Error(t, err)
	assert.Equal(t, tier.ErrNotFound, err)
}

func TestPromoteCustomerTier(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store)

	merchantId := uuid.New()
	customerId := uuid.New()
	//
	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "white",
		MinPoints:  100,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "red",
		MinPoints:  2000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "black",
		MinPoints:  1000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	err := tierSvc.ProcessOrder(ctx, customerId, merchantId, 200)
	assert.NoError(t, err)
	err = tierSvc.PromoteCustomerTier(ctx, merchantId, customerId)
	assert.NoError(t, err)
	customerAcc, err := tierSvc.GetTierCustomer(ctx, merchantId, customerId)
	assert.NoError(t, err)
	assert.Equal(t, "black", customerAcc.TierName)
	//
	tierSvc.PromoteCustomerTier(ctx, merchantId, customerId)
	//
	err = tierSvc.PromoteCustomerTier(ctx, merchantId, customerId)
	assert.Error(t, err)
	assert.Equal(t, tier.ErrHighTier, err)
}

func TestDemoteCustomerTier(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	tierSvc := tier.NewTierSvc(store, store)

	merchantId := uuid.New()
	customerId := uuid.New()
	//
	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "white",
		MinPoints:  100,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "red",
		MinPoints:  2000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	tierSvc.CreateTierLevel(ctx, tier.TierLevel{
		ID:         uuid.New(),
		MerchantID: merchantId,
		Name:       "black",
		MinPoints:  1000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	//
	err := tierSvc.ProcessOrder(ctx, customerId, merchantId, 1100)
	assert.NoError(t, err)
	err = tierSvc.DemoteCustomerTier(ctx, merchantId, customerId)
	assert.NoError(t, err)
	customerAcc, err := tierSvc.GetTierCustomer(ctx, merchantId, customerId)
	assert.NoError(t, err)
	assert.Equal(t, "white", customerAcc.TierName)
	//
	err = tierSvc.DemoteCustomerTier(ctx, merchantId, customerId)
	assert.Error(t, err)
	assert.Equal(t, tier.ErrLowTier, err)
}
