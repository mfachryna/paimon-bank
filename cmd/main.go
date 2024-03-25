package main

import (
	"github.com/mfachryna/paimon-bank/config"
	"github.com/mfachryna/paimon-bank/internal/app"
)

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
