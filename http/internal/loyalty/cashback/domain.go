package cashback

import (
	"time"

	"github.com/google/uuid"
)

type Cashback struct {
	ID               uuid.UUID
	MerchantID       uuid.UUID
	Name             string
	Description      string
	Percentage       float64
	EligibleProducts []uuid.UUID
	Active           bool
	Expiration       time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
type CashbackUser struct {
	MerchantID uuid.UUID
	CashbackID uuid.UUID
	UserID     uuid.UUID
	Points     float64
}
