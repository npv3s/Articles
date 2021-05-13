package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"main/database"
	"net/http"
	"strconv"
)

func (h *Handler) Articles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.database.GetArticles()
	if err != nil {
		log.Println("Get all articles error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	tp, err := template.New("").ParseFiles("templates/base.html", "templates/articles.html")
	if err != nil {
		log.Fatal("Template rendering error:", err)
	}

	err = tp.ExecuteTemplate(w, "base", struct {
		Title   string
		Content interface{}
	}{
		"Статьи",
		articles,
	})
	if err != nil {
		log.Println("Template rendering error:", err)
	}
}

func (h *Handler) Article(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Unknown article id:", vars["id"])
	}

	//user := h.GetUser(r)

	article, err := h.database.GetArticle(articleId)
	if err != nil {
		log.Println("Get", articleId, "article error:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	tp, err := template.New("").ParseFiles("templates/base.html", "templates/article.html", "templates/comment.html", "templates/reply.html")
	if err != nil {
		log.Fatal("Template rendering error:", err)
	}

	rootComment := struct {
		Id   int
		Root *int
	}{
		0,
		nil,
	}

	content := struct {
		Article      *database.Article
		RootComment  interface{}
		IsOwner      bool
		IsAuthorized bool
	}{
		article,
		//&article.Author == user,
		rootComment,
		false,
		false,
	}

	err = tp.ExecuteTemplate(w, "base", struct {
		Title   string
		Content interface{}
	}{
		article.Title,
		content,
	})
	if err != nil {
		log.Println("Template rendering error:", err)
	}
}

func (h *Handler) ArticleUpdate(w http.ResponseWriter, r *http.Request) {
	type ArticleUpdate struct {
		Id    int    `json:"article_id"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	var articleUpdate ArticleUpdate

	err := json.NewDecoder(r.Body).Decode(&articleUpdate)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.UpdateArticle(1, articleUpdate.Id, articleUpdate.Title, articleUpdate.Body)
	if err != nil {
		log.Println("Article update error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ArticleDelete(w http.ResponseWriter, r *http.Request) {
	type ArticleDelete struct {
		Id int `json:"article_id"`
	}

	var articleDelete ArticleDelete

	err := json.NewDecoder(r.Body).Decode(&articleDelete)
	if err != nil {
		log.Println("Json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.database.DeleteArticle(1, articleDelete.Id)
	if err != nil {
		log.Println("Article delete error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
