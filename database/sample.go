package database

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type SampleDB struct{}

func NewSample() SampleDB {
	return SampleDB{}
}

func (_ *SampleDB) NewUser(login, password string) (*int, error) {
	fmt.Println("Login:", login)
	fmt.Println("password:", password)
	id := 0
	return &id, nil
}

func (_ *SampleDB) GetUserByLogin(login string) (*User, error) {
	if login == "npv3s" {
		return &User{
			1,
			"npv3s",
			"12345",
		}, nil
	} else {
		return nil, errors.New("no such user")
	}
}

func (_ *SampleDB) GetUserById(userId int) (*User, error) {
	if userId == 1 {
		return &User{
			1,
			"npv3s",
			"12345",
		}, nil
	} else {
		return nil, errors.New("no such user")
	}
}

func (_ *SampleDB) GetArticles() ([]ArticleDescription, error) {
	return []ArticleDescription{
		{0, "npv3s", "Hi", time.Now()},
		{1, "npv3s", "Bye", time.Now()},
	}, nil
}

func (_ *SampleDB) GetArticle(id int) (*Article, error) {
	text := "В течение последних 20 лет сотрудники Post Office (почтовая компания из Великобритании) " +
		"разбирались с программой Horizon, в которой имелась фатальная ошибка: из-за неисправности казалось, " +
		"что сотрудники воровали десятки тысяч фунтов. Некоторые местные почтмейстеры были осуждены " +
		"и посажены в тюрьму из-за того, что Post Office упорно настаивал на том, что программному обеспечению " +
		"можно доверять. После десятилетий баталий приговоры 39 человек, наконец, отменили. Случай стал крупнейшей " +
		"судебной ошибкой, которую когда-либо видела Великобритания."
	digits := []int{0, 1, 2, 3, 4, 5}
	return &Article{
		ArticleDescription{id, "npv3s", "Hi", time.Now()}, 1, text, []Comment{
			{nil, 1, "npv3s", "Hello", time.Now(), nil},
			{nil, 2, "npv3s", "Bye", time.Now(), []Comment{
				{&digits[1], 3, "abc", "Goodbye", time.Now(), []Comment{}},
			}}},
	}, nil
}

func (_ *SampleDB) UpdateArticle(authorId, articleId int, title, body string) error {
	log.Println("Article update:", title)
	return nil
}

func (_ *SampleDB) DeleteArticle(authorId, articleId int) error {
	return nil
}

func (_ *SampleDB) GetComments(articleId int) ([]Comment, error) {
	one := 1
	return []Comment{
		{nil, 1, "npv3s", "Hello", time.Now(), nil},
		{nil, 2, "npv3s", "Bye", time.Now(), []Comment{
			{&one, 3, "abc", "Goodbye", time.Now(), []Comment{}},
		}},
	}, nil
}

func (_ *SampleDB) NewComment(authorId, articleId int, text string, root *int) error {
	fmt.Println("New comment:", authorId, text)
	return nil
}

func (_ *SampleDB) UpdateComment(authorId, commentId int, text string) error {
	fmt.Println("Update comment:", commentId, text)
	return nil
}

func (_ *SampleDB) DeleteComment(authorId, commentId int) error {
	fmt.Println("Delete comment:", commentId)
	return nil
}

func (_ *SampleDB) NewArticle(text string) (*int, error) {
	one := 1
	return &one, nil
}