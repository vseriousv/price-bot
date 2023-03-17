package models

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	Id                  int64  `json:"id"`
	ChatId              int64  `json:"chat_id"`
	UserName            string `json:"user_name"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Description         string `json:"description"`
	Photo               string `json:"photo"`
	Title               string `json:"title"`
	AllMembersAreAdmins bool   `json:"all_members_are_admins"`
	InviteLink          string `json:"invite_link"`
}

func (u *User) Create(db *pgxpool.Pool) error {
	query := `
INSERT INTO users(chat_id, user_name, first_name, last_name, description, photo, title, all_members_are_admins, invite_link, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now(), now());
`
	_, err := db.Exec(context.Background(), query,
		u.ChatId,
		u.UserName,
		u.FirstName,
		u.LastName,
		u.Description,
		u.Photo,
		u.Title,
		u.AllMembersAreAdmins,
		u.InviteLink,
	)
	if err != nil {
		return err
	}
	return nil
}
