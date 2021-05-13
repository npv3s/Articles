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
	Body     string
	Comments []Comment
}

type Comment struct {
	Root     *int
	Id       int
	Author   string
	Text     string
	Time     time.Time
	Comments []Comment
}

type Database interface {
	NewUser(login, password string) (*int, error)
	GetUserByLogin(login string) (*User, error)
	GetUserById(userId int) (*User, error)

	GetArticles() (*[]ArticleDescription, error)
	GetArticle(id int) (*Article, error)
	UpdateArticle(authorId, articleId int, title, body string) error
	DeleteArticle(authorId, articleId int) error

	GetComments(articleId int) ([]Comment, error)
	NewComment(authorId, articleId int, text string, root *int) error
	UpdateComment(authorId, commentId int, text string) error
	DeleteComment(authorId, commentId int) error
}
