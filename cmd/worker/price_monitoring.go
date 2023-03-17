package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/price-bot/internal/app"
	"github.com/vseriousv/price-bot/internal/config"
	"github.com/vseriousv/price-bot/pkg/utils"
	"log"
)

const serviceName = "PRICE_MONITORING_WORKER"

func main() {
	c := config.DefaultConfig()

	if err := run(c); err != nil {
		log.Fatal("[ERROR]", err)
	}
}

func run(c *config.Config) error {
	// connect database
	pgxConf, err := pgxpool.ParseConfig(c.DbUrl)
	utils.HandleError(err, "DB is not connection")
	pgxConf.MinConns = 1
	pgxConf.MaxConns = 8

	dbContext := context.Background()
	pool, err := pgxpool.ConnectConfig(dbContext, pgxConf)
	utils.HandleError(err, "DB is not connection")

	err = pool.Ping(context.Background())
	utils.HandleError(err, "DB is not connection")

	// create service
	service := app.New(serviceName, pool, c)
	err = app.NewPriceMonitoringWorkerService(service).Run()
	utils.HandleError(err, "Failed to start PriceMonitoringWorker service")
	return nil
}
