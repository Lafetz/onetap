package tier

import (
	"time"

	"github.com/google/uuid"
)

type TierLevel struct {
	ID         uuid.UUID
	MerchantID uuid.UUID
	Name       string
	MinPoints  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type CustomerTier struct {
	MerchantID uuid.UUID
	CustomerID uuid.UUID
	TierName   string
	Points     int
}

func NewCustomerTier(MerchantID uuid.UUID, customerID uuid.UUID, tierName string, points int) CustomerTier {
	return CustomerTier{
		MerchantID: MerchantID,
		CustomerID: customerID,
		TierName:   tierName,
		Points:     points,
	}
}
func NewTierLevel(MerchantID uuid.UUID, name string, minPoints int) TierLevel {
	return TierLevel{
		MerchantID: MerchantID,
		Name:       name,
		MinPoints:  minPoints,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
