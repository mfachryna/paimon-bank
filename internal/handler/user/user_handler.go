package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mfachryna/paimon-bank/config"
	interfaces "github.com/mfachryna/paimon-bank/internal/interfaces"
	"github.com/mfachryna/paimon-bank/pkg/jwt"
	"go.uber.org/zap"
)

type UserHandler struct {
	ur  interfaces.UserRepository
	val *validator.Validate
	cfg config.Configuration
	log *zap.Logger
}

func NewUserHandler(
	r chi.Router,
	ur interfaces.UserRepository,
	val *validator.Validate,
	cfg config.Configuration,
	log *zap.Logger,
) {
	uh := &UserHandler{
		ur:  ur,
		val: val,
		cfg: cfg,
		log: log,
	}

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", uh.Register)
		r.Post("/login", uh.Login)

		r.Route("/", func(r chi.Router) {
			r.Use(jwt.JwtMiddleware)
		})
	})
}
