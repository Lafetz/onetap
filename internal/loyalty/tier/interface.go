package tier

import (
	"context"

	"github.com/google/uuid"
)

type TierLevelRepository interface {
	CreateTierLevel(ctx context.Context, tierLevel TierLevel) error
	GetTierLevel(ctx context.Context, merchantId uuid.UUID, name string) (TierLevel, error)
	UpdateTierLevel(ctx context.Context, updatedTierLevel TierLevel) error
	DeleteTierLevel(ctx context.Context, merchantId uuid.UUID, name string) error
	ListTierLevel(ctx context.Context, merchantId uuid.UUID) ([]TierLevel, error)
}
type CustomerTierRepository interface {
	CreateCustomerTier(ctx context.Context, customerTier CustomerTier) error
	GetCustomerTier(ctx context.Context, merchantId uuid.UUID, customerId uuid.UUID) (CustomerTier, error)
	UpdateCustomerTier(ctx context.Context, updatedCustomerTier CustomerTier) error
	DeleteCustomerTier(ctx context.Context, merchantId uuid.UUID, customerId uuid.UUID) error
}
