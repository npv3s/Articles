package database

import "time"

type User struct {
	Id       int
	Login    string
	PassHash string
}

type Tag string

type Article struct {
	Id     int
	Author string
	Title  string
	Body   string
	Time   time.Time
	Tags   []Tag
}

type Comment struct {
	Id     int
	Root   *int
	Author string
	Text   string
	Time   time.Time
}

type Database interface {
	NewUser(login, password string) error
	GetPassword(login string) (*User, error)

	GetArticles() (*[]Article, error)
	GetArticle(id int) (*Article, error)
	UpdateArticle(id int, body string) error
	DeleteArticle(id int) error

	GetComments(articleId int) ([]Comment, error)
	NewComment(comment Comment) error
	UpdateComment(comment Comment) error
	DeleteComment(commentId int) error
}
