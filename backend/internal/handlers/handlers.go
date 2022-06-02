package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Sitri-code/observer_bot/server/internal/models"
	"github.com/Sitri-code/observer_bot/server/internal/repository"
	"net/http"
	"strconv"
	"time"
)

func errorHandler(err error, msg string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(msg, err)))
}

type Handlers struct {
	r *repository.Repository
}

func New(r *repository.Repository) *Handlers {
	return &Handlers{r}
}

func (h *Handlers) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}
	msgs, err := h.r.GetMessages()
	if err != nil {
		errorHandler(err, "Failed to get all data: %v", w)
		return
	}
	enc, err := json.Marshal(msgs)
	if err != nil {
		errorHandler(err, "Failed to get all data: %v", w)
		return
	}
	w.Write([]byte(enc))
}

func (h *Handlers) GetMessageByTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	themeId := r.FormValue("id")

	val, err := strconv.Atoi(themeId)
	if err != nil {
		errorHandler(err, "Failed to parse data: %v", w)
		return
	}

	m, err := h.r.GetMessagesByThemeId(val)
	if err != nil {
		errorHandler(err, "Failed to get data: %v", w)
		return
	}
	enc, err := json.Marshal(m)
	if err != nil {
		errorHandler(err, "Failed to get all data: %v", w)
		return
	}
	w.Write([]byte(enc))
}

func (h *Handlers) GetUsersByTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	themeId := r.FormValue("theme_id")

	val, err := strconv.Atoi(themeId)
	if err != nil {
		errorHandler(err, "Failed to parse data: %v", w)
		return
	}

	m, err := h.r.GetUsersByThemeId(val)
	if err != nil {
		errorHandler(err, "Failed to get data: %v", w)
		return
	}
	enc, err := json.Marshal(m)
	if err != nil {
		errorHandler(err, "Failed to get all data: %v", w)
		return
	}
	w.Write([]byte(enc))
}

func (h *Handlers) GetUsersByChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	chatId := r.FormValue("chat_id")

	val, err := strconv.Atoi(chatId)
	if err != nil {
		errorHandler(err, "Failed to parse data: %v", w)
		return
	}

	m, err := h.r.GetUsersByChatId(val)
	if err != nil {
		errorHandler(err, "Failed to get data: %v", w)
		return
	}
	enc, err := json.Marshal(m)
	if err != nil {
		errorHandler(err, "Failed to get all data: %v", w)
		return
	}
	w.Write([]byte(enc))
}

func (h *Handlers) PostMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	id := r.FormValue("id")
	chatId := r.FormValue("chat_id")
	themeId := r.FormValue("theme_id")
	message := r.FormValue("message")
	timeNow := time.Now().String()
	ref := r.FormValue("ref")

	chatVal, chatErr := strconv.Atoi(chatId)
	idVal, idErr := strconv.Atoi(id)
	themeVal, themeErr := strconv.Atoi(themeId)

	if chatErr != nil || idErr != nil || themeErr != nil || message == "" || ref == "" {
		errorHandler(chatErr, "Failed to parse data: %v", w)
		return
	}

	msg := models.Message{Id: idVal, ChatId: chatVal, ThemeId: themeVal, Message: message, Time: timeNow, Ref: ref}
	_, err := h.r.CreateMessage(msg)
	if err != nil {
		errorHandler(err, "Failed to post message: %v", w)
		return
	}
	w.Write([]byte("Success create message"))
}

func (h *Handlers) PostUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	chatId := r.FormValue("chat_id")
	name := r.FormValue("name")
	reference := r.FormValue("reference")
	themeId := r.FormValue("theme_id")

	chatVal, chatErr := strconv.Atoi(chatId)
	themeVal, themeErr := strconv.Atoi(themeId)
	if chatErr != nil || themeErr != nil || name == "" || reference == "" || themeId == "" || chatId == "" {
		errorHandler(chatErr, "Failed to parse data: %v", w)
		return
	}
	user := models.User{ChatId: chatVal, Name: name, Reference: reference, ThemeId: themeVal}
	_, err := h.r.CreateUser(user)
	if err != nil {
		errorHandler(err, "Failed to post user: %v", w)
		return
	}
	w.Write([]byte("Success create user"))
}

func (h *Handlers) PostTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	id := r.FormValue("theme_id")
	val, err := strconv.Atoi(id)
	if err != nil {
		errorHandler(err, "Failed to parse data: %v", w)
		return
	}
	name := r.FormValue("name")

	theme := models.Theme{Id: val, Name: name}
	_, err = h.r.CreateTheme(theme)
	if err != nil {
		errorHandler(err, "Failed to post theme: %v", w)
		return
	}
	w.Write([]byte("Success create theme"))
}

func (h *Handlers) DeleteThemeById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}

	id := r.FormValue("id")

	val, err := strconv.Atoi(id)
	if err != nil {
		errorHandler(err, "Failed to parse data: %v", w)
		return
	}

	err = h.r.DeleteTheme(val)
	if err != nil {
		errorHandler(err, "Failed to delete theme: %v", w)
		return
	}
	w.Write([]byte("Success delete theme"))
}

func (h *Handlers) PutUserTheme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}
	themeId := r.FormValue("theme_id")
	chatId := r.FormValue("chat_id")

	themeVar, themeErr := strconv.Atoi(themeId)
	if themeErr != nil {
		errorHandler(themeErr, "Failed to parse data: %v", w)
		return
	}

	chatVar, chatErr := strconv.Atoi(chatId)
	if chatErr != nil {
		errorHandler(chatErr, "Failed to parse data: %v", w)
		return
	}
	err := h.r.UpdateUser(themeVar, chatVar)
	if err != nil {
		errorHandler(err, "Failed to update theme: %v", w)
		return
	}
	w.Write([]byte("Success update user theme"))
}

func (h *Handlers) GetThemeIdByChatId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("method %v not allowed", r.Method)))
		return
	}
	chatId := r.FormValue("chat_id")
	res, err := h.r.GetThemeId(chatId)
	if err != nil {
		errorHandler(err, "Failed to get data: %v", w)
	}
	w.Write([]byte(res))
}
