package db

import (
	"context"
	"fmt"
	"log"
	"time"

	pgxpool_prometheus "github.com/cmackenzie1/pgxpool-prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mfachryna/paimon-bank/config"
	"github.com/prometheus/client_golang/prometheus"
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

	if err := PingDatabase(context.Background(), dbPool); err != nil {
		log.Fatalf("Can't pinging database: %v", err)
	}

	prometheus.MustRegister(pgxpool_prometheus.NewPgxPoolStatsCollector(dbPool, c.Postgres.PostgresqlDbname))

	log.Println("Success connect to database")
	return dbPool
}

func PingDatabase(ctx context.Context, pgx *pgxpool.Pool) error {
	if err := pgx.Ping(ctx); err != nil {
		return err
	}

	return nil
}
