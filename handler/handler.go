package handler

import (
	"main/authenticator"
	"main/database"
	"net/http"
	"strconv"
)

type baseTemplate struct {
	Title string
	IsAuthorized bool
	Content interface{}
}

type Handler struct {
	authenticator authenticator.Authenticator
	database      interface{ database.Database }
}

func NewHandler(database interface{ database.Database }) Handler {
	return Handler{
		authenticator.NewAuthenticator(database),
		database,
	}
}

func (h *Handler) GetUserId(r *http.Request) *int {
	userIdCookie, err := r.Cookie("user_id")
	if err != nil {
		return nil
	}

	userId, err := strconv.Atoi(userIdCookie.Value)
	if err != nil {
		return nil
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}

	if h.authenticator.CheckSession(userId, tokenCookie.Value) {
		return &userId
	}

	return nil
}
