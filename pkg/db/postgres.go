package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mfachryna/paimon-bank/config"
)

// Return new Postgresql db instance
func NewPsqlDB(c *config.Configuration) *pgxpool.Pool {
	dataSourceName := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s",
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlPassword,
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresParams,
	)

	log.Println("db conn", dataSourceName)
	poolConf, err := pgxpool.ParseConfig(dataSourceName)
	if err != nil {
		log.Fatalf("Error when parsing db config: %v", err)
	}

	poolConf.MaxConns = 100
	poolConf.MaxConnLifetime = time.Hour
	poolConf.MaxConnIdleTime = time.Minute * 30
	poolConf.ConnConfig.ConnectTimeout = time.Second * 5

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		log.Fatalf("Error when creating db pool: %v", err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Can't pinging database: %v", err)
	}

	log.Println("Successed connect to databases")
	return dbPool
}
