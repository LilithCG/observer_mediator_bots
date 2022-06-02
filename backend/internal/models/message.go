package models

type Message struct {
	Id      int    `json:"id"`
	ChatId  int    `json:"chat_id"`
	ThemeId int    `json:"theme_id"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Ref     string `json:"ref"`
}
