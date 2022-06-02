package models

type User struct {
	ChatId    int    `json:"chat_id"`
	Name      string `json:"name"`
	Reference string `json:"ref"`
	ThemeId   int    `json:"theme_id"`
}
