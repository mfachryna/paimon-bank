package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mfachryna/paimon-bank/config"
	interfaces "github.com/mfachryna/paimon-bank/internal/interfaces"
	"github.com/mfachryna/paimon-bank/pkg/jwt"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	ur  interfaces.UserRepository
	br  interfaces.BalanceRepository
	val *validator.Validate
	cfg config.Configuration
	log *zap.Logger
}

func NewBalanceHandler(
	r chi.Router,
	ur interfaces.UserRepository,
	br interfaces.BalanceRepository,
	val *validator.Validate,
	cfg config.Configuration,
	log *zap.Logger,
) {
	bh := &BalanceHandler{
		ur:  ur,
		br:  br,
		val: val,
		cfg: cfg,
		log: log,
	}
	r.Group(func(r chi.Router) {
		r.Use(jwt.JwtMiddleware)
		r.Post("/transaction", bh.Transaction)
		r.Route("/balance", func(r chi.Router) {
			r.Post("/", bh.Store)
			r.Get("/", bh.Index)
			r.Get("/history", bh.History)
		})
	})
}
