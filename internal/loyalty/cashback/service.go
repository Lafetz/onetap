package cashback

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrNotAuth           = errors.New("unauth")
	ErrNotActiveCashback = errors.New("merchant doesn't have active cashback")
	ErrActiveCashback    = errors.New("merchant can only have one active cashback")
)

type CashbackSvc struct {
	cashbackRepo     CashBackRepository
	cashbackUserRepo CashbackCustRepository
}

func NewCashbackSvc(cashbackRepo CashBackRepository, cashbackUserRepo CashbackCustRepository) *CashbackSvc {
	return &CashbackSvc{
		cashbackRepo:     cashbackRepo,
		cashbackUserRepo: cashbackUserRepo,
	}
}
