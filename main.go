package main

import (
	"fmt"
	"log"
	"main/database"
	"main/handler"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	db, err := database.NewPgPool(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(db)

	r := mux.NewRouter()

	r.Path("/login/").Methods("GET").HandlerFunc(h.LoginForm)

	r.Path("/login/").Methods("POST").HandlerFunc(h.Login)

	r.Path("/signup/").Methods("GET").HandlerFunc(h.SignUpForm)

	r.Path("/signup/").Methods("POST").HandlerFunc(h.SignUp)

	r.Path("/article/new/").Methods("GET").HandlerFunc(h.ArticleForm)

	r.Path("/article/new/").Methods("POST").HandlerFunc(h.ArticleNew)

	r.Path("/article/{id:[0-9]+}/").HandlerFunc(h.Article)

	r.Path("/article/update/").HandlerFunc(h.ArticleUpdate)

	r.Path("/article/delete/").HandlerFunc(h.ArticleDelete)

	r.Path("/comment/new/").HandlerFunc(h.AddComment)

	r.PathPrefix("/front/").Handler(http.StripPrefix("/front/", http.FileServer(http.Dir("front/"))))

	r.Path("/").HandlerFunc(h.Articles)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
