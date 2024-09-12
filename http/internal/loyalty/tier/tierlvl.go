package tier

import (
	"context"

	"github.com/google/uuid"
)

func (t *TierSvc) CreateTierLevel(ctx context.Context, tierLevel TierLevel) error {
	err := t.tierLevelRepo.CreateTierLevel(ctx, tierLevel)
	if err != nil {
		return err
	}

	return nil
}
func (t *TierSvc) ListTierLevels(ctx context.Context, merchantId uuid.UUID) ([]TierLevel, error) {
	tierLevels, err := t.tierLevelRepo.ListTierLevel(ctx, merchantId)
	if err != nil {
		return nil, err
	}
	return tierLevels, nil
}

func (t *TierSvc) GetTierLevel(ctx context.Context, merchantId uuid.UUID, name string) (TierLevel, error) {
	tierLevel, err := t.tierLevelRepo.GetTierLevel(ctx, merchantId, name)
	if err != nil {
		return TierLevel{}, err
	}
	return tierLevel, nil
}

func (t *TierSvc) DeleteTierLevel(ctx context.Context, merchantId uuid.UUID, name string) error {
	err := t.tierLevelRepo.DeleteTierLevel(ctx, merchantId, name)
	if err != nil {
		return err
	}

	return nil
}
func (t *TierSvc) UpdateTierLevel(ctx context.Context, updatedTierLevel TierLevel) error {
	err := t.tierLevelRepo.UpdateTierLevel(ctx, updatedTierLevel)
	if err != nil {
		return err
	}
	return nil
}
