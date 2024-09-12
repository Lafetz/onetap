package integration_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/loyalty/cashback"
	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateCashbackUser(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	newCashbackUser := cashback.CashbackUser{
		MerchantID: uuid.New(),
		CashbackID: uuid.New(),
		UserID:     uuid.New(),
		Points:     150.0,
	}

	err := cashbackSvc.CreateCashbackUser(ctx, newCashbackUser)
	assert.NoError(t, err)
}

func TestGetCashbackUser(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	newCashbackUser := cashback.CashbackUser{
		MerchantID: uuid.New(),
		CashbackID: uuid.New(),
		UserID:     uuid.New(),
		Points:     150.0,
	}

	err := cashbackSvc.CreateCashbackUser(ctx, newCashbackUser)
	assert.NoError(t, err)

	cashbackUser, err := cashbackSvc.GetCashbackUser(ctx, newCashbackUser.UserID, newCashbackUser.CashbackID)
	assert.NoError(t, err)
	assert.Equal(t, newCashbackUser.CashbackID, cashbackUser.CashbackID)
	assert.Equal(t, newCashbackUser.UserID, cashbackUser.UserID)
}

func TestUpdateCashbackUser(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	newCashbackUser := cashback.CashbackUser{
		MerchantID: uuid.New(),
		CashbackID: uuid.New(),
		UserID:     uuid.New(),
		Points:     150.0,
	}

	err := cashbackSvc.CreateCashbackUser(ctx, newCashbackUser)
	assert.NoError(t, err)

	updatedCashbackUser := newCashbackUser
	updatedCashbackUser.Points = 200.0
	err = cashbackSvc.UpdateCashbackUser(ctx, updatedCashbackUser)
	assert.NoError(t, err)

	cashbackUser, err := cashbackSvc.GetCashbackUser(ctx, newCashbackUser.UserID, newCashbackUser.CashbackID)
	assert.NoError(t, err)
	assert.Equal(t, updatedCashbackUser.Points, cashbackUser.Points)
}

func TestDeleteCashbackUser(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	newCashbackUser := cashback.CashbackUser{
		MerchantID: uuid.New(),
		CashbackID: uuid.New(),
		UserID:     uuid.New(),
		Points:     150.0,
	}

	err := cashbackSvc.CreateCashbackUser(ctx, newCashbackUser)
	assert.NoError(t, err)

	err = cashbackSvc.DeleteCashbackUser(ctx, newCashbackUser.UserID, newCashbackUser.CashbackID)
	assert.NoError(t, err)

	_, err = cashbackSvc.GetCashbackUser(ctx, newCashbackUser.UserID, newCashbackUser.CashbackID)
	assert.Error(t, err)
	assert.Equal(t, cashback.ErrNotFound, err)
}

func TestListCashbackUsers(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	merchantID := uuid.New()

	newCashbackUser := cashback.CashbackUser{
		MerchantID: merchantID,
		CashbackID: uuid.New(),
		UserID:     uuid.New(),
		Points:     150.0,
	}

	err := cashbackSvc.CreateCashbackUser(ctx, newCashbackUser)
	assert.NoError(t, err)

	cashbackUsers, err := cashbackSvc.ListCashbackUsers(ctx, merchantID)
	assert.NoError(t, err)
	assert.Len(t, cashbackUsers, 1)
	assert.Equal(t, newCashbackUser.UserID, cashbackUsers[0].UserID)
}

func TestProcessOrder_NoPreviousSell(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	merchantID := uuid.New()
	userID := uuid.New()
	productId := uuid.New()

	cashbackToCreate := cashback.Cashback{
		ID:               uuid.New(),
		MerchantID:       merchantID,
		Name:             "First Sale Cashback",
		Description:      "Cashback for first-time sales",
		Percentage:       0.05, // 5% cashback
		EligibleProducts: []uuid.UUID{productId},
		Active:           true,
		Expiration:       time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := cashbackSvc.CreateCashback(ctx, cashbackToCreate)
	assert.NoError(t, err)

	err = cashbackSvc.ProcessOrder(ctx, userID, merchantID, productId, 200)
	assert.NoError(t, err)

	cashbackUser, err := cashbackSvc.GetCashbackUser(ctx, userID, cashbackToCreate.ID)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, cashbackUser.Points) // 5% of 200 = 10 points
}

func TestProcessOrder_NewCashback(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	merchantID := uuid.New()
	userID := uuid.New()
	productId := uuid.New()

	newCashback := cashback.Cashback{
		ID:               uuid.New(),
		MerchantID:       merchantID,
		Name:             "New Cashback Offer",
		Description:      "Discount for new purchases",
		Percentage:       0.1,
		EligibleProducts: []uuid.UUID{productId, uuid.New(), uuid.New()},
		Active:           true,
		Expiration:       time.Now().Add(60 * 24 * time.Hour),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := cashbackSvc.CreateCashback(ctx, newCashback)
	assert.NoError(t, err)

	err = cashbackSvc.ProcessOrder(ctx, userID, merchantID, productId, 300.0)
	assert.NoError(t, err)
	err = cashbackSvc.ProcessOrder(ctx, userID, merchantID, productId, 300.0)
	assert.NoError(t, err)

	cashbackUser, err := cashbackSvc.GetCashbackUser(ctx, userID, newCashback.ID)
	assert.NoError(t, err)
	assert.Equal(t, 60.0, cashbackUser.Points)
}

func TestProcessOrder_ExpiredCashback(t *testing.T) {
	ctx := context.Background()
	store := repository.NewDb(testDbInstance)
	cashbackSvc := cashback.NewCashbackSvc(store, store, &mockNot{}, slog.Default())

	merchantID := uuid.New()
	userID := uuid.New()
	productId := uuid.New()

	expiredCashback := cashback.Cashback{
		ID:               uuid.New(),
		MerchantID:       merchantID,
		Name:             "Expired Cashback",
		Description:      "Expired cashback offer",
		Percentage:       0.2,
		EligibleProducts: []uuid.UUID{productId},
		Active:           false,
		Expiration:       time.Now().Add(-5 * 24 * time.Hour), // Expired 5 days ago
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := cashbackSvc.CreateCashback(ctx, expiredCashback)
	assert.NoError(t, err)

	err = cashbackSvc.ProcessOrder(ctx, userID, merchantID, productId, 300.0)
	assert.NoError(t, err)
	_, err = cashbackSvc.GetCashbackUser(ctx, userID, expiredCashback.ID)
	assert.Error(t, err)

}
