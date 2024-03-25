package interfaces

import (
	"context"

	"github.com/mfachryna/paimon-bank/internal/entity"
)

// Translation -.
type (
	UserRepository interface {
		Get(context.Context, entity.User) error
		FindById(context.Context, string) (*entity.User, error)
		FindByEmail(context.Context, string) (*entity.User, error)
		FindByPhone(context.Context, string) (*entity.User, error)
		Insert(context.Context, entity.User) error
		Delete(context.Context, string) error
		Update(context.Context, entity.User) error
		EmailCheck(context.Context, string) (int64, error)
	}
)
