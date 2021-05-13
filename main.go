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

	db := database.NewSample()
	h := handler.NewHandler(&db)

	r := mux.NewRouter()

	r.Path("/article/{id:[0-9]+}/").HandlerFunc(h.Article)

	r.Path("/article/update/{id:[0-9]+}/").HandlerFunc(h.ArticleUpdate)

	r.Path("/article/delete/{id:[0-9]+}/").HandlerFunc(h.ArticleDelete)

	r.Path("/comment/new/{id:[0-9]+}/").HandlerFunc(h.AddComment)

	r.PathPrefix("/front/").Handler(http.StripPrefix("/front/", http.FileServer(http.Dir("front/"))))

	r.Path("/").HandlerFunc(h.Articles)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
