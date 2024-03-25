package interfaces

import (
	"context"

	dto "github.com/mfachryna/paimon-bank/internal/domain/dto/balance"
	"github.com/mfachryna/paimon-bank/internal/entity"
)

// Translation -.
type (
	BalanceRepository interface {
		GetBalance(context.Context, string, string) (int64, error)
		GetBalanceList(context.Context, dto.BalanceFilter, string) ([]entity.Balance, int64, error)
		GetBalanceHistory(context.Context, dto.BalanceFilter, string) ([]entity.BalanceHistory, int64, error)
		Insert(context.Context, entity.BalanceHistory) error
	}
)
