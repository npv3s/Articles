package handler

import (
	"html/template"
	"log"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	tp, err := template.New("").ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		log.Fatal("Template rendering error:", err)
	}

	content := struct {}{}

	err = tp.ExecuteTemplate(w, "base", struct {
		Title   string
		Content interface{}
	}{
		"Вход",
		content,
	})

	if err != nil {
		http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
