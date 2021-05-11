package handler

import (
	"main/authenticator"
	"main/database"
	"net/http"
)

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

func (h *Handler) GetUser(r *http.Request) *string {
	loginCookie, err := r.Cookie("login")
	if err != nil {
		return nil
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return nil
	}

	if h.authenticator.CheckSession(loginCookie.Value, tokenCookie.Value) {
		return &loginCookie.Value
	}

	return nil
}
