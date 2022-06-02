package repository

import (
	"database/sql"
	"github.com/Sitri-code/observer_bot/server/internal/models"
	"github.com/Sitri-code/observer_bot/server/pkg/db"
)

type Repository struct {
	db *sql.DB
}

func New() (*Repository, error) {
	d, err := db.Init()
	if err != nil {
		return nil, err
	}
	return &Repository{d}, nil
}

func (r *Repository) CreateMessage(m models.Message) (int64, error) {
	query, err := r.db.Prepare("INSERT INTO message (id, chat_id, theme_id, message, time, ref) VALUES ($1, $2, $3, $4, $5, $6);")
	if err != nil {
		return 0, err
	}

	res, err := query.Exec(m.Id, m.ChatId, m.ThemeId, m.Message, m.Time, m.Ref)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	defer query.Close()
	return count, nil
}

func (r *Repository) CreateUser(u models.User) (int64, error) {
	query, err := r.db.Prepare("INSERT INTO users (chat_id, name, ref, theme_id) VALUES ($1, $2, $3, $4);")
	if err != nil {
		return 0, err
	}

	res, err := query.Exec(u.ChatId, u.Name, u.Reference, u.ThemeId)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	defer query.Close()
	return count, nil
}

func (r *Repository) CreateTheme(t models.Theme) (int64, error) {
	query, err := r.db.Prepare("INSERT INTO theme (theme_id, name) VALUES ($1, $2);")
	if err != nil {
		return 0, err
	}

	res, err := query.Exec(t.Id, t.Name)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	defer query.Close()
	return count, nil
}

func (r *Repository) GetMessages() ([]models.Message, error) {
	rows, err := r.db.Query("SELECT * FROM message")
	if err != nil {
		return nil, err
	}

	var msgs []models.Message
	for rows.Next() {
		var id int
		var chatId int
		var themeId int
		var message string
		var Time string
		var ref string
		err := rows.Scan(&id, &chatId, &themeId, &message, &Time, &ref)
		if err != nil {
			return nil, err
		}
		msg := models.Message{Id: id, ChatId: chatId, ThemeId: themeId, Message: message, Time: Time, Ref: ref}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (r *Repository) GetMessagesByThemeId(themeId int) ([]models.Message, error) {
	rows, err := r.db.Query(`SELECT * FROM message WHERE theme_id = $1`, themeId)
	if err != nil {
		return nil, err
	}

	var msgs []models.Message
	for rows.Next() {
		var id int
		var chatId int
		var themeId int
		var message string
		var Time string
		var ref string
		err := rows.Scan(&id, &chatId, &themeId, &message, &Time, &ref)
		if err != nil {
			return nil, err
		}
		msg := models.Message{Id: id, ChatId: chatId, ThemeId: themeId, Message: message, Time: Time, Ref: ref}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (r *Repository) GetUsersByThemeId(themeId int) ([]models.User, error) {
	rows, err := r.db.Query(`SELECT * FROM users WHERE theme_id = $1`, themeId)
	if err != nil {
		return nil, err
	}

	var msgs []models.User
	for rows.Next() {
		var chatId int
		var name string
		var ref string
		var theme int
		err := rows.Scan(&chatId, &name, &ref, &theme)
		if err != nil {
			return nil, err
		}
		msg := models.User{ChatId: chatId, Name: name, Reference: ref, ThemeId: theme}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (r *Repository) GetUsersByChatId(chatId int) ([]models.User, error) {
	rows, err := r.db.Query(`SELECT * FROM users WHERE chat_id = $1`, chatId)
	if err != nil {
		return nil, err
	}

	var msgs []models.User
	for rows.Next() {
		var chat int
		var name string
		var ref string
		var themeId int
		err := rows.Scan(&chat, &name, &ref, &themeId)
		if err != nil {
			return nil, err
		}
		msg := models.User{ChatId: chat, Name: name, Reference: ref, ThemeId: themeId}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (r *Repository) DeleteTheme(themeId int) error {
	_, err := r.db.Query("DELETE FROM theme WHERE theme_id = $1", themeId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateUser(themeId int, chatId int) error {
	_, err := r.db.Query("UPDATE users SET theme_id = ? WHERE chat_id = $1", themeId, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetThemeId(chatId string) (string, error) {
	var id string
	err := r.db.QueryRow(`SELECT theme_id FROM users WHERE chat_id = $1`, chatId).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
