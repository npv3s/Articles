package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	type NewComment struct {
		ArticleId int    `json:"article_id"`
		Root      *int   `json:"root"`
		Text      string `json:"body"`
	}

	var comment NewComment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.NewComment(1, comment.ArticleId, comment.Text, comment.Root)
	if err != nil {
		log.Println("Comment add error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	type NewComment struct {
		ArticleId int    `json:"article_id"`
		Root      *int   `json:"root"`
		Text      string `json:"body"`
	}

	var comment NewComment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.UpdateComment(1, comment.ArticleId, comment.Text)
	if err != nil {
		log.Println("Comment add error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	type NewComment struct {
		ArticleId int    `json:"article_id"`
		Root      *int   `json:"root"`
		Text      string `json:"body"`
	}

	var comment NewComment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.NewComment(1, comment.ArticleId, comment.Text, comment.Root)
	if err != nil {
		log.Println("Comment add error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
