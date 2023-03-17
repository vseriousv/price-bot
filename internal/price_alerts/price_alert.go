package price_alerts

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/price-bot/internal/models"
	"log"
)

func GetById(db *pgxpool.Pool, chatID int64, id int64) (*models.PriceAlert, error) {
	query := `
SELECT pa.id,
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
       (extract(epoch from pa.created_at) * 1000)::bigint as created_at,
	   pa.create_price < pa.alert_price as is_up
FROM price_alerts pa
LEFT JOIN users u on u.id = pa.user_id
WHERE u.chat_id = $1 AND pa.id = $2;
`
	var p models.PriceAlert
	err := db.QueryRow(context.Background(), query, chatID, id).
		Scan(
			&p.Id,
			&p.User,
			&p.Ticker,
			&p.CreatePrice,
			&p.AlertPrice,
			&p.CreateAt,
			&p.IsUp,
		)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetListByChatId(db *pgxpool.Pool, chatID int64) (*[]models.PriceAlert, error) {
	query := `
SELECT pa.id,
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
       (extract(epoch from pa.created_at) * 1000)::bigint as created_at,
	   pa.create_price < pa.alert_price as is_up
FROM price_alerts pa
LEFT JOIN users u on u.id = pa.user_id
WHERE u.chat_id = $1
ORDER BY pa.created_at DESC;
`
	rows, err := db.Query(context.Background(), query, chatID)
	if err != nil {
		return nil, err
	}

	var list []models.PriceAlert
	for rows.Next() {
		var p models.PriceAlert
		err := rows.Scan(
			&p.Id,
			&p.User,
			&p.Ticker,
			&p.CreatePrice,
			&p.AlertPrice,
			&p.CreateAt,
			&p.IsUp,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, p)
	}

	return &list, nil
}

func DeleteAlertById(db *pgxpool.Pool, id, userId int64) (bool, error) {
	query := `
delete from price_alerts pa
where pa.id = $1 and pa.user_id = $2;
`
	_, err := db.Exec(context.Background(), query, id, userId)
	if err != nil {
		log.Println("Opps, ", err)
	}
	return false, nil
}
