package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Unknown article id:", vars["id"])
	}

	type NewComment struct {
		ArticleId int `json:"article_id"`
		Root *int `json:"root"`
		Text string `json:"body"`
	}

	var comment NewComment

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.NewComment(articleId, "npv3s", comment.Text, comment.Root)
	if err != nil {
		log.Println("Comment add error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
