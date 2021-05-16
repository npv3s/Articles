package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	err = h.database.NewComment(*user, comment.ArticleId, comment.Text, comment.Root)
	if err != nil {
		log.Println("Comment add error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	err = h.database.UpdateComment(*user, comment.ArticleId, comment.Text)
	if err != nil {
		log.Println("Comment add error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type DeleteComment struct {
		CommentId int    `json:"comment_id"`
	}

	var comment DeleteComment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.DeleteComment(*user, comment.CommentId)
	if err != nil {
		log.Println("Comment delete error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
