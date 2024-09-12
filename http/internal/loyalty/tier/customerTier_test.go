package tier

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCustomerTier(t *testing.T) {
	// Define tier levels
	tierLevels := []TierLevel{
		{
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Name:       "Bronze",
			MinPoints:  50,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Name:       "Silver",
			MinPoints:  1000,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			MerchantID: uuid.New(),
			Name:       "Gold",
			MinPoints:  1500,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	customerAcc := CustomerTier{
		CustomerID: uuid.New(),
		MerchantID: uuid.New(),
		TierName:   "Bronze",
		Points:     40,
	}

	updatedCustomerAcc := UpdateCustomerTier(customerAcc, 10, tierLevels)

	assert.Equal(t, 50, updatedCustomerAcc.Points)
	assert.Equal(t, "Bronze", updatedCustomerAcc.TierName)

	customerAcc = CustomerTier{
		CustomerID: uuid.New(),
		MerchantID: uuid.New(),
		TierName:   "Bronze",
		Points:     900,
	}

	updatedCustomerAcc = UpdateCustomerTier(customerAcc, 150, tierLevels)

	assert.Equal(t, 1050, updatedCustomerAcc.Points)
	assert.Equal(t, "Silver", updatedCustomerAcc.TierName)

	customerAcc = CustomerTier{
		CustomerID: uuid.New(),
		MerchantID: uuid.New(),
		TierName:   "Silver",
		Points:     950,
	}

	updatedCustomerAcc = UpdateCustomerTier(customerAcc, 50, tierLevels)

	assert.Equal(t, 1000, updatedCustomerAcc.Points)
	assert.Equal(t, "Silver", updatedCustomerAcc.TierName)
}
func TestPromoteTier(t *testing.T) {
	tests := []struct {
		name         string
		customerTier CustomerTier
		tierLevels   []TierLevel
		expectedTier string
		expectedErr  error
	}{
		{
			name: "Successful promotion",
			customerTier: CustomerTier{
				TierName: "Bronze",
			},
			tierLevels: []TierLevel{
				{Name: "Bronze", MinPoints: 50},
				{Name: "Silver", MinPoints: 1000},
				{Name: "Gold", MinPoints: 1500},
			},
			expectedTier: "Silver",
			expectedErr:  nil,
		},
		{
			name: "No tier to promote",
			customerTier: CustomerTier{
				TierName: "Gold",
			},
			tierLevels: []TierLevel{
				{Name: "Bronze", MinPoints: 50},
				{Name: "Silver", MinPoints: 1000},
				{Name: "Gold", MinPoints: 1500},
			},
			expectedTier: "Gold",
			expectedErr:  ErrHighTier,
		},
		{
			name: "Tier not found",
			customerTier: CustomerTier{
				TierName: "NonExistent",
			},
			tierLevels: []TierLevel{
				{Name: "Bronze", MinPoints: 50},
				{Name: "Silver", MinPoints: 1000},
			},
			expectedTier: "NonExistent",
			expectedErr:  ErrNoTier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PromoteTier(tt.customerTier, tt.tierLevels)

			if err != nil && err != tt.expectedErr {
				t.Errorf("PromoteTier() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if result.TierName != tt.expectedTier {
				t.Errorf("PromoteTier() got = %v, want %v", result.TierName, tt.expectedTier)
			}
		})
	}
}

func TestDemoteTier(t *testing.T) {
	tests := []struct {
		name         string
		customerTier CustomerTier
		tierLevels   []TierLevel
		expectedTier string
		expectedErr  error
	}{
		{
			name: "Successful demotion",
			customerTier: CustomerTier{
				TierName: "Silver",
			},
			tierLevels: []TierLevel{
				{Name: "Bronze", MinPoints: 50},
				{Name: "Silver", MinPoints: 1000},
				{Name: "Gold", MinPoints: 1500},
			},
			expectedTier: "Bronze",
			expectedErr:  nil,
		},
		{
			name: "No tier to demote",
			customerTier: CustomerTier{
				TierName: "Bronze",
			},
			tierLevels: []TierLevel{
				{Name: "Bronze", MinPoints: 50},
			},
			expectedTier: "Bronze",
			expectedErr:  ErrLowTier,
		},
		{
			name: "Tier not found",
			customerTier: CustomerTier{
				TierName: "NonExistent",
			},
			tierLevels: []TierLevel{
				{Name: "Bronze", MinPoints: 50},
				{Name: "Silver", MinPoints: 1000},
			},
			expectedTier: "NonExistent",
			expectedErr:  ErrNoTier,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DemoteTier(tt.customerTier, tt.tierLevels)

			if err != nil && err != tt.expectedErr {
				t.Errorf("DemoteTier() error = %v, wantErr %v", err, tt.expectedErr)
				return
			}
			if result.TierName != tt.expectedTier {
				t.Errorf("DemoteTier() got = %v, want %v", result.TierName, tt.expectedTier)
			}
		})
	}
}
