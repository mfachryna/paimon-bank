package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mfachryna/paimon-bank/internal/entity"
	"go.uber.org/zap"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewUserRepo(db *pgxpool.Pool, log *zap.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (ur *UserRepository) Get(ctx context.Context, data entity.User) error {
	return nil
}

func (ur *UserRepository) FindById(ctx context.Context, userId string) (*entity.User, error) {
	res := &entity.User{}
	sql := `SELECT id, name, email, password FROM users WHERE id = $1`

	err := ur.db.QueryRow(ctx, sql, userId).Scan(&res.ID, &res.Name, &res.Email, &res.Password)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	res := &entity.User{}
	sql := `SELECT id, name, email, password FROM users WHERE email = $1`

	err := ur.db.QueryRow(ctx, sql, email).Scan(&res.ID, &res.Name, &res.Email, &res.Password)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ur *UserRepository) Insert(ctx context.Context, data entity.User) error {
	var sql string

	sql = `INSERT INTO users (id,email,name,password) VALUES ($1,$2,$3,$4)`
	if _, err := ur.db.Exec(ctx, sql, data.ID, data.Email, data.Name, data.Password); err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, userId string) error {
	return nil
}

func (ur *UserRepository) Update(ctx context.Context, data entity.User) error {
	sql := `UPDATE users SET name = $1, email = $2, password = $3, WHERE id = $4`

	_, err := ur.db.Exec(ctx, sql, data.Name, data.Email, data.Password, data.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) EmailCheck(ctx context.Context, email string) (int64, error) {
	var count int64

	if err := ur.db.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE email = $1", email).Scan(&count); err != nil {
		return 0, nil
	}

	return count, nil
}

func (ur *UserRepository) PhoneCheck(ctx context.Context, phone string) (int64, error) {
	var count int64

	if err := ur.db.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE phone = $1", phone).Scan(&count); err != nil {
		return 0, nil
	}

	return count, nil
}
