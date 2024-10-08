// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package gen

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

type CustomerTier struct {
	MerchantID uuid.UUID
	CustomerID uuid.UUID
	TierName   string
	Points     int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TierLevel struct {
	TierID     uuid.UUID
	MerchantID uuid.UUID
	Name       string
	MinPoints  int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
