package models

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PriceAlert struct {
	Id                 int64   `json:"id"`
	User               User    `json:"user"`
	Ticker             string  `json:"ticker"`
	CreatePrice        float64 `json:"create_price"`
	AlertPrice         float64 `json:"alert_price"`
	CreateAt           int64   `json:"create_at"`
	IsUp               bool    `json:"is_up"`
	ObservablePercent  float64 `json:"observable_percent"`
	SubscriptionTicker string  `json:"subscription_ticker"`
}

func (u *PriceAlert) Create(db *pgxpool.Pool) error {
	query := `
INSERT INTO price_alerts(user_id, ticker, create_price, alert_price, created_at)
VALUES ($1, $2, $3, $4, now());
`
	_, err := db.Exec(context.Background(), query,
		u.User.Id,
		u.Ticker,
		u.CreatePrice,
		u.AlertPrice,
	)
	if err != nil {
		return err
	}
	return nil
}
