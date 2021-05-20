package database

import "time"

type User struct {
	Id       int
	Login    string
	PassHash string
}

type ArticleDescription struct {
	Id     int
	Author string
	Title  string
	Time   time.Time
}

type Article struct {
	ArticleDescription
	AuthorId int
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
	NewUser(login, passHash string) (*int, error)
	GetUserByLogin(login string) (*User, error)
	GetUserById(userId int) (*User, error)
	GetUsers() ([]User, error)
	GetAdmins() ([]int, error)
	NewPassword(userId int, passHash string) error
	IsAdmin(userId int) (bool, error)

	NewArticle(authorId int, title, text string) (*int, error)
	GetArticles() ([]ArticleDescription, error)
	GetArticle(id int) (*Article, error)
	UpdateArticle(userId, articleId int, title, body string) error
	DeleteArticle(userId, articleId int) error

	GetComments(articleId int) ([]Comment, error)
	NewComment(authorId, articleId int, text string, root *int) error
	UpdateComment(authorId, commentId int, text string) error
	DeleteComment(authorId, commentId int) error
}
