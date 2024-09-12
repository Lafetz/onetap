package cashback

import (
	"errors"
	"log/slog"
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
	notification     notification
	logger           *slog.Logger
}

func NewCashbackSvc(cashbackRepo CashBackRepository, cashbackUserRepo CashbackCustRepository, notification notification, logger *slog.Logger) *CashbackSvc {
	return &CashbackSvc{
		cashbackRepo:     cashbackRepo,
		cashbackUserRepo: cashbackUserRepo,
		logger:           logger,
		notification:     notification,
	}
}
