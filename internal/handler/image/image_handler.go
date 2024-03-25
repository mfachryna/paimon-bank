package image

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mfachryna/paimon-bank/config"
	"github.com/mfachryna/paimon-bank/pkg/jwt"
	"go.uber.org/zap"
)

type ImageHandler struct {
	val *validator.Validate
	cfg config.Configuration
	log *zap.Logger
}

func NewImageHandler(r chi.Router, val validator.Validate, cfg config.Configuration, log *zap.Logger) {
	ih := &ImageHandler{
		val: &val,
		cfg: cfg,
		log: log,
	}

	r.Route("/image", func(r chi.Router) {
		r.Use(jwt.JwtMiddleware)
		r.Post("/", ih.Store)
	})
}
