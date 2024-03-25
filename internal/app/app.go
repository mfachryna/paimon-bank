package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/mfachryna/paimon-bank/config"
	"github.com/mfachryna/paimon-bank/internal/common/utils/validation"
	"github.com/mfachryna/paimon-bank/pkg/db"
	"github.com/mfachryna/paimon-bank/pkg/logger"
	"github.com/mfachryna/paimon-bank/pkg/promotheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(cfg *config.Configuration) {
	var validate *validator.Validate

	pgx := db.NewPsqlDB(cfg)

	validate = validator.New()
	if err := validation.RegisterCustomValidation(validate); err != nil {
		log.Fatalf("error register custom validation")
	}

	logger, err := logger.Initialize(*cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())
	r.Route("/v1", func(r chi.Router) {
		r.Use(promotheus.PrometheusMiddleware)
	})

	s := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}
	go func() {
		fmt.Println("Listen and Serve at port 8080")
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("error in ListenAndServe: %s", err)
		}
	}()
	log.Print("Server Started")

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt)
	<-stopped

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("shutting down gracefully...")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("error in Server Shutdown: %s", err)
	}
	fmt.Println("server stopped")
}
