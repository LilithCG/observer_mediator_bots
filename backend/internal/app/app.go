package app

import (
	"github.com/Sitri-code/observer_bot/server/internal/handlers"
	"github.com/Sitri-code/observer_bot/server/internal/repository"
	"net/http"
)

func Run(port string) error {
	repo, err := repository.New()
	if err != nil {
		return err
	}
	hand := handlers.New(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/getall", hand.GetAllMessages)
	mux.HandleFunc("/api/get", hand.GetMessageByTheme)
	mux.HandleFunc("/api/createmsg", hand.PostMessage)
	mux.HandleFunc("/api/createuser", hand.PostUser)
	mux.HandleFunc("/api/createtheme", hand.PostTheme)
	mux.HandleFunc("/api/deletetheme", hand.DeleteThemeById)
	mux.HandleFunc("/api/update", hand.PutUserTheme)
	mux.HandleFunc("/api/get_theme_id", hand.GetThemeIdByChatId)
	mux.HandleFunc("/api/get_users_theme", hand.GetUsersByTheme)
	mux.HandleFunc("/api/get_users_chat", hand.GetUsersByChat)

	return http.ListenAndServe(":"+port, mux)
}
