package tier

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func UpdateCustomerTier(customerAcc CustomerTier, amount int, tierLevels []TierLevel) CustomerTier {
	customerAcc.Points += amount
	var currentTierLevel TierLevel
	for _, tierLevel := range tierLevels {
		if tierLevel.Name == customerAcc.TierName {
			currentTierLevel = tierLevel
			break
		}
	}
	if currentTierLevel.MinPoints <= customerAcc.Points {
		for _, tierLevel := range tierLevels {
			if tierLevel.MinPoints > currentTierLevel.MinPoints && customerAcc.Points >= tierLevel.MinPoints {
				customerAcc.TierName = tierLevel.Name
			}
		}
	}
	return customerAcc
}
func (t *TierSvc) ProcessOrder(ctx context.Context, customerId, merchantId uuid.UUID, amount int) error { // gets triggered each time user orders from merchant
	tierLevels, err := t.tierLevelRepo.ListTierLevel(ctx, merchantId)
	if err != nil {

		return err
	}
	if len(tierLevels) == 0 { //merchant has no tiers
		return nil
	}
	customerAcc, err := t.custoTierRepo.GetCustomerTier(ctx, merchantId, customerId)
	if err != nil {
		if errors.Is(err, ErrNotFound) {

			customerAcc = NewCustomerTier(merchantId, customerId, tierLevels[0].Name, 0) //create new customer with 0 points
			customerAcc = UpdateCustomerTier(customerAcc, amount, tierLevels)
			return t.custoTierRepo.CreateCustomerTier(ctx, customerAcc)
		}

		return err
	}

	customerAcc = UpdateCustomerTier(customerAcc, amount, tierLevels)
	return t.custoTierRepo.UpdateCustomerTier(ctx, customerAcc)
}

func (t *TierSvc) GetTierCustomer(ctx context.Context, merchantId, customerId uuid.UUID) (CustomerTier, error) {
	customerTier, err := t.custoTierRepo.GetCustomerTier(ctx, merchantId, customerId)
	if err != nil {
		return CustomerTier{}, err
	}
	if customerTier.MerchantID != merchantId || customerTier.CustomerID != customerId {
		return CustomerTier{}, ErrUnauth
	}
	return customerTier, nil
}
func PromoteTier(customerTier CustomerTier, tierLevels []TierLevel) (CustomerTier, error) {
	var currentTierIndex int = -1
	for i, tierLevel := range tierLevels {
		if tierLevel.Name == customerTier.TierName {
			currentTierIndex = i
			break
		}
	}
	if currentTierIndex == -1 {
		return customerTier, ErrNoTier
	}
	if currentTierIndex+1 < len(tierLevels) {
		customerTier.TierName = tierLevels[currentTierIndex+1].Name
		return customerTier, nil
	}
	return customerTier, ErrHighTier
}
func (t *TierSvc) PromoteCustomerTier(ctx context.Context, merchantId, customerId uuid.UUID) error {
	customerTier, err := t.custoTierRepo.GetCustomerTier(ctx, merchantId, customerId)
	if err != nil {
		return err
	}

	tierLevels, err := t.tierLevelRepo.ListTierLevel(ctx, merchantId)
	if err != nil {
		return err
	}

	promotedCustomerTier, err := PromoteTier(customerTier, tierLevels)
	if err != nil {
		return err
	}

	return t.custoTierRepo.UpdateCustomerTier(ctx, promotedCustomerTier)
}
func DemoteTier(customerTier CustomerTier, tierLevels []TierLevel) (CustomerTier, error) {
	var currentTierIndex int = -1
	for i, tierLevel := range tierLevels {
		if tierLevel.Name == customerTier.TierName {
			currentTierIndex = i
			break
		}
	}
	if currentTierIndex == -1 {
		return customerTier, ErrNoTier
	}

	if currentTierIndex-1 >= 0 {
		customerTier.TierName = tierLevels[currentTierIndex-1].Name
		return customerTier, nil
	}

	return customerTier, ErrLowTier
}
func (t *TierSvc) DemoteCustomerTier(ctx context.Context, merchantId, customerId uuid.UUID) error {
	customerTier, err := t.custoTierRepo.GetCustomerTier(ctx, merchantId, customerId)
	if err != nil {
		return err
	}

	tierLevels, err := t.tierLevelRepo.ListTierLevel(ctx, merchantId)
	if err != nil {
		return err
	}

	demotedCustomerTier, err := DemoteTier(customerTier, tierLevels)
	if err != nil {
		return err
	}

	return t.custoTierRepo.UpdateCustomerTier(ctx, demotedCustomerTier)
}
