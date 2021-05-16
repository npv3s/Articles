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

func (h *Handler) Articles(w http.ResponseWriter, _ *http.Request) {
	articles, err := h.database.GetArticles()
	if err != nil {
		http.Error(w, "Get all articles error: "+err.Error(), http.StatusInternalServerError)
		return
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

	user := h.GetUserId(r)

	article, err := h.database.GetArticle(articleId)
	if err != nil {
		http.Error(w, "Get "+strconv.Itoa(articleId)+", article error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	isOwner := false
	if user != nil {
		if *user == article.AuthorId {
			isOwner = true
		}
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
		rootComment,
		isOwner,
		user != nil,
	}

	err = tp.ExecuteTemplate(w, "base", struct {
		Title   string
		Content interface{}
	}{
		article.Title,
		content,
	})
	if err != nil {
		http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ArticleUpdate(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	err = h.database.UpdateArticle(*user, articleUpdate.Id, articleUpdate.Title, articleUpdate.Body)
	if err != nil {
		log.Println("Article update error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ArticleDelete(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	err = h.database.DeleteArticle(*user, articleDelete.Id)
	if err != nil {
		http.Error(w,"Article delete error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ArticleForm(w http.ResponseWriter, _ *http.Request) {
	tp, err := template.New("").ParseFiles("templates/base.html", "templates/new-article.html")
	if err != nil {
		log.Fatal("Template rendering error:", err)
	}

	content := struct{}{}

	err = tp.ExecuteTemplate(w, "base", struct {
		Title   string
		Content interface{}
	}{
		"Новая статья",
		content,
	})

	if err != nil {
		http.Error(w, "Template rendering error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ArticleNew(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type NewArticle struct {
		Title string `json:"title"`
		Text string `json:"body"`
	}

	var article NewArticle

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Println("json parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	articleId, err := h.database.NewArticle(*user, article.Title, article.Text)
	if err != nil {
		http.Error(w, "article add error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct {
		ArticleId int `json:"article_id"`
	}{
		*articleId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
