package users

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vseriousv/price-bot/internal/models"
)

func GetByChatId(db *pgxpool.Pool, chatID int64) (*models.User, error) {
	var u models.User
	query := `
SELECT id, chat_id, user_name, first_name, last_name, description, photo, title, all_members_are_admins, invite_link FROM users
WHERE chat_id = $1;
`
	err := db.QueryRow(context.Background(), query, chatID).
		Scan(
			&u.Id,
			&u.ChatId,
			&u.UserName,
			&u.FirstName,
			&u.LastName,
			&u.Description,
			&u.Photo,
			&u.Title,
			&u.AllMembersAreAdmins,
			&u.InviteLink,
		)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
