package price_monitoring_worker

import (
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/price-bot/internal/models"
	"github.com/vseriousv/price-bot/internal/providers"
	"log"
	"sync"
	"time"
)

func checkObservableAlert(wg *sync.WaitGroup, db *pgxpool.Pool, token string) {
	taskRun := func() {
		tickers := getTickers(db)

		p, err := providers.GetProvider("kucoin")
		if err != nil {
			log.Println(err)
		}

		var executedAlerts []models.PriceAlert
		for _, ticker := range tickers {
			price := string(*p.GetPriceByTicker(ticker))
			item := getObservableAlerts(db, ticker, price)
			executedAlerts = append(executedAlerts, item...)
		}
		sendAlertObservable(token, executedAlerts)
	}

	taskRun()
	for range time.Tick(time.Second * 30) {
		taskRun()
	}
	wg.Done()
}

func getObservableAlerts(db *pgxpool.Pool, ticker, price string) []models.PriceAlert {
	queryTickers := `
select pa.id,
       json_build_object(
               'id', u.id,
               'chat_id', u.chat_id,
               'user_name', u.user_name,
               'first_name', u.first_name,
               'last_name', u.last_name,
               'description', u.description,
               'photo', u.photo,
               'title', u.title,
               'all_members_are_admins', u.all_members_are_admins,
               'invite_link', u.invite_link
           ) as user,
       pa.ticker,
       pa.create_price,
       pa.alert_price,
       (extract(epoch from pa.created_at) * 1000)::bigint as created_at
from price_alerts pa
   left join users u on u.id = pa.user_id
where (
              (pa.create_price > pa.alert_price and pa.alert_price > $1)
              or (pa.create_price < pa.alert_price and pa.alert_price < $1)
          )
and pa.ticker = $2;
`
	rows, err := db.Query(context.Background(), queryTickers, price, ticker)
	if err != nil {
		log.Println(err)
		return nil
	}

	var res []models.PriceAlert
	for rows.Next() {
		var pa models.PriceAlert
		err := rows.Scan(
			&pa.Id,
			&pa.User,
			&pa.Ticker,
			&pa.CreatePrice,
			&pa.AlertPrice,
			&pa.CreateAt,
		)
		if err != nil {
			log.Println("getExecutedAlerts", err)
			return nil
		}
		res = append(res, pa)
	}

	return res
}

func sendAlertObservable(token string, alerts []models.PriceAlert) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	var alertIds []int64
	for _, alert := range alerts {
		text := fmt.Sprintf("[EXECUTED ALERT] :: [%s] :: %.8g", alert.Ticker, alert.AlertPrice)
		msg := tgbotapi.NewMessage(alert.User.ChatId, text)
		_, err := bot.Send(msg)
		if err == nil {
			alertIds = append(alertIds, alert.Id)
		}
	}
}
