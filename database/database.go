package database

import "time"

type User struct {
	Id       int
	Login    string
	PassHash string
}

type Tag string

type ArticleDescription struct {
	Id     int
	Author string
	Title  string
	Time   time.Time
	Tags   []Tag
}

type Article struct {
	ArticleDescription
	Body   string
	Comments []Comment
}

type Comment struct {
	Id     int
	Author string
	Text   string
	Time   time.Time
	Comments []Comment
}

type Database interface {
	NewUser(login, password string) error
	GetPassword(login string) (*User, error)

	GetArticles() (*[]ArticleDescription, error)
	GetArticle(id int) (*Article, error)
	UpdateArticle(id int, title, body string) error
	DeleteArticle(id int) error

	GetComments(articleId int) ([]Comment, error)
	NewComment(articleId int, author, text string, root *int) error
	UpdateComment(comment Comment) error
	DeleteComment(commentId int) error
}
