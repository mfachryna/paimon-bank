package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	dto "github.com/mfachryna/paimon-bank/internal/domain/dto/balance"
	"github.com/mfachryna/paimon-bank/internal/entity"
	"go.uber.org/zap"
)

type BalanceRepository struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewBalanceRepo(db *pgxpool.Pool, log *zap.Logger) *BalanceRepository {
	return &BalanceRepository{
		db:  db,
		log: log,
	}
}

func (br *BalanceRepository) GetBalanceHistory(ctx context.Context, filter dto.BalanceFilter, userId string) ([]entity.BalanceHistory, int64, error) {
	sql := fmt.Sprintf(`SELECT 
	id,
	balance,
	bank_number,
	bank_name,
	currency,
	receipt,
	user_id,
	created_at
	FROM balance_histories
	WHERE user_id = '%s' 
	ORDER BY created_at desc 
	LIMIT %d OFFSET %d`, userId, filter.Limit, filter.Offset)

	rows, err := br.db.Query(ctx, sql)
	if err != nil {
		return []entity.BalanceHistory{}, 0, err
	}

	data := make([]entity.BalanceHistory, 0)
	var count int64 = 0
	for rows.Next() {
		var createdAt time.Time
		var balanceHistory entity.BalanceHistory
		err := rows.Scan(
			&balanceHistory.ID,
			&balanceHistory.Balance,
			&balanceHistory.Source.BankNumber,
			&balanceHistory.Source.BankName,
			&balanceHistory.Currency,
			&balanceHistory.Receipt,
			&balanceHistory.UserId,
			&createdAt,
		)

		balanceHistory.CreatedAt = fmt.Sprint(createdAt.Unix())
		if err != nil {
			return []entity.BalanceHistory{}, 0, err
		}
		data = append(data, balanceHistory)
		count += 1
	}
	rows.Close()

	return data, count, nil
}

func (br *BalanceRepository) GetBalanceList(ctx context.Context, filter dto.BalanceFilter, userId string) ([]entity.Balance, int64, error) {
	sql := fmt.Sprintf(`SELECT balance,currency FROM balances WHERE user_id = '%s' ORDER BY ID desc LIMIT %d OFFSET %d`, userId, filter.Limit, filter.Offset)
	rows, err := br.db.Query(ctx, sql)
	if err != nil {
		return []entity.Balance{}, 0, err
	}

	data := make([]entity.Balance, 0)
	var count int64 = 0
	for rows.Next() {
		var balance entity.Balance
		err := rows.Scan(
			&balance.Balance,
			&balance.Currency,
		)
		if err != nil {
			return []entity.Balance{}, 0, err
		}
		data = append(data, balance)
		count += 1
	}
	rows.Close()

	return data, count, nil
}

func (br *BalanceRepository) GetBalance(ctx context.Context, userId string, currency string) (int64, error) {
	var sql string

	var balance int64
	sql = `SELECT sum(balance) FROM balances WHERE user_id = $1 and currency = $2`
	if err := br.db.QueryRow(ctx, sql, userId, currency).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}

func (br *BalanceRepository) Insert(ctx context.Context, data entity.BalanceHistory) error {
	var sql string

	sql = `INSERT INTO balance_histories (id,balance,bank_number,bank_name,currency,receipt,user_id,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,now())`
	if _, err := br.db.Exec(ctx, sql, data.ID, data.Balance, data.Source.BankNumber, data.Source.BankName, data.Currency, data.Receipt, data.UserId); err != nil {
		return err
	}

	sql = `INSERT INTO balances (balance,currency,user_id) VALUES ($1,$2,$3) ON CONFLICT (currency,user_id) DO UPDATE SET balance = balances.balance + EXCLUDED.balance`
	if _, err := br.db.Exec(ctx, sql, data.Balance, data.Currency, data.UserId); err != nil {
		return err
	}

	return nil
}
