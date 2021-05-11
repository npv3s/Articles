package database

import (
	"errors"
	"fmt"
	"time"
)

type SampleDB struct{}

func NewSample() SampleDB {
	return SampleDB{}
}

func (_ *SampleDB) NewUser(login, password string) error {
	fmt.Println("Login:", login)
	fmt.Println("password:", password)
	return nil
}

func (_ *SampleDB) GetPassword(login string) (*User, error) {
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

func (_ *SampleDB) GetArticles() (*[]Article, error) {
	return &[]Article{
		{0, "npv3s", "Hi", "", time.Now(), nil},
		{1, "npv3s", "Bye", "", time.Now(), nil},
	}, nil
}

func (_ *SampleDB) GetArticle(id int) (*Article, error) {
	text := "В течение последних 20 лет сотрудники Post Office (почтовая компания из Великобритании) " +
		"разбирались с программой Horizon, в которой имелась фатальная ошибка: из-за неисправности казалось, " +
		"что сотрудники воровали десятки тысяч фунтов. Некоторые местные почтмейстеры были осуждены " +
		"и посажены в тюрьму из-за того, что Post Office упорно настаивал на том, что программному обеспечению " +
		"можно доверять. После десятилетий баталий приговоры 39 человек, наконец, отменили. Случай стал крупнейшей " +
		"судебной ошибкой, которую когда-либо видела Великобритания."
	return &Article{
		id, "npv3s", "Hi", text, time.Now(), nil,
	}, nil
}

func (_ *SampleDB) UpdateArticle(id int, body string) error {
	return nil
}

func (_ *SampleDB) DeleteArticle(id int) error {
	return nil
}

func (_ *SampleDB) GetComments(articleId int) ([]Comment, error) {
	one := 1
	return []Comment{
		{1, nil, "npv3s", "Hello", time.Now()},
		{2, &one, "npv3s", "Bye", time.Now()},
	}, nil
}

func (_ *SampleDB) NewComment(comment Comment) error {
	return nil
}

func (_ *SampleDB) UpdateComment(comment Comment) error {
	return nil
}

func (_ *SampleDB) DeleteComment(commentId int) error {
	return nil
}
