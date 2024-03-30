package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mfachryna/paimon-bank/internal/entity"
	"github.com/mfachryna/paimon-bank/pkg/db"
)

type HeathHandler struct {
	pgx *pgxpool.Pool
}

func NewHealthHandler(r chi.Router, pgx *pgxpool.Pool) {
	hh := &HeathHandler{
		pgx: pgx,
	}

	r.Get("/healthz", hh.Healthz)
}

func (hh *HeathHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	health := entity.HealthResponse{
		Service: "ok",
	}

	if err := db.PingDatabase(context.Background(), hh.pgx); err != nil {
		health.Service = "database error"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
