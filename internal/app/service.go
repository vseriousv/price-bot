package app

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/price-bot/internal/config"
)

type service struct {
	Name   string
	Db     *pgxpool.Pool
	Config *config.Config
}

func New(name string, db *pgxpool.Pool, c *config.Config) *service {
	return &service{
		Name:   name,
		Db:     db,
		Config: c,
	}
}
