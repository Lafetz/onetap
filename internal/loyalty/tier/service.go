package tier

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrMerchantNoTier = errors.New("merchant has no tiers")
	ErrUnauth         = errors.New("Unauthorized")
	ErrDepulicateTier = errors.New("another tier with the same name exists")
	ErrNoTier         = errors.New("customer doesn't have any tier")
	ErrLowTier        = errors.New("cusomer already on lowest tier")
	ErrHighTier       = errors.New("already at the highest tier")
)

type TierSvc struct {
	tierLevelRepo TierLevelRepository
	custoTierRepo CustomerTierRepository
}

func NewTierSvc(tierLevelRepo TierLevelRepository, custoTierRepo CustomerTierRepository) *TierSvc {
	return &TierSvc{
		tierLevelRepo: tierLevelRepo,
		custoTierRepo: custoTierRepo,
	}
}
