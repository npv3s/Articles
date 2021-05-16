package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PgDB struct {
	pool *pgxpool.Pool
}

func NewPgPool(url string) (*PgDB, error) {
	dbPool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	createSQL := `BEGIN;


CREATE TABLE IF NOT EXISTS public."user"
(
    id integer NOT NULL,
    login character varying(250),
    pass_hash character(60),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.article
(
    id integer NOT NULL,
    author integer NOT NULL,
    title character varying(250) NOT NULL,
    body text NOT NULL,
    created timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.comment
(
    id integer NOT NULL,
    root integer,
    author integer NOT NULL,
    comment character varying(500) NOT NULL,
    article integer NOT NULL,
    created timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE public.article
    ADD FOREIGN KEY (author)
    REFERENCES public."user" (id)
    NOT VALID;


ALTER TABLE public.comment
    ADD FOREIGN KEY (author)
    REFERENCES public."user" (id)
    NOT VALID;


ALTER TABLE public.comment
    ADD FOREIGN KEY (article)
    REFERENCES public.article (id)
    NOT VALID;

END;`

	_, err = dbPool.Exec(context.Background(), createSQL)
	if err != nil {
		return nil, err
	}

	return &PgDB{
		dbPool,
	}, nil
}

func (db *PgDB) NewUser(login, passHash string) (*int, error) {
	r, err := db.pool.Exec(context.Background(), "INSERT INTO user VALUES(?, ?)", login, passHash)
	if err != nil {
		return nil, err
	} else if r.RowsAffected() == 0 {
		return nil, errors.New("zero affected rows")
	}

	var userID int
	err = db.pool.QueryRow(context.Background(), "SELECT id FROM user WHERE login = ?", login).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (db *PgDB) GetUserByLogin(login string) (*User, error) {
	panic("implement me")
}

func (db *PgDB) GetUserById(userId int) (*User, error) {
	panic("implement me")
}

func (db *PgDB) NewArticle(text string) (*int, error) {
	panic("implement me")
}

func (db *PgDB) GetArticles() ([]ArticleDescription, error) {
	sql := `SELECT id, (SELECT login FROM public."user" WHERE "user".id = author), title, created FROM public.article`
	rows, err := db.pool.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}

	var articles []ArticleDescription

	defer rows.Close()
	for rows.Next() {
		var id int
		var author, title string
		var created time.Time
		err = rows.Scan(&id, &author, &title, &created)
		if err != nil {
			return nil, err
		}
		articles = append(articles, ArticleDescription{
			id,
			author,
			title,
			created,
		})
	}

	return articles, err
}

func (db *PgDB) GetArticle(id int) (*Article, error) {
	sql := `SELECT author, (SELECT login FROM public."user" WHERE "user".id = author), title, body, created FROM public.article WHERE id = $1`
	var authorId int
	var author, title, body string
	var created time.Time
	err := db.pool.QueryRow(context.Background(), sql, id).Scan(&authorId, &author, &title, &body, &created)
	if err != nil {
		return nil, err
	}

	comments, err := db.GetComments(id)
	if err != nil {
		return nil, err
	}

	return &Article{
		ArticleDescription{
			id, author, title, created,
		},
		authorId,
		body,
		comments,
	}, nil
}

func (db *PgDB) UpdateArticle(authorId, articleId int, title, body string) error {
	panic("implement me")
}

func (db *PgDB) DeleteArticle(authorId, articleId int) error {
	panic("implement me")
}

func commentsR(root *int, comments []Comment) []Comment {
	var parsedComments []Comment

	for _, comment := range comments {
		if comment.Root == root {
			parsedComments = append(parsedComments, Comment{
				comment.Root,
				comment.Id,
				comment.Author,
				comment.Text,
				comment.Time,
				commentsR(&comment.Id, comments),
			})
		}
	}

	return parsedComments
}

func (db *PgDB) GetComments(articleId int) ([]Comment, error) {
	sql := `SELECT id, root, (SELECT login FROM public."user" WHERE "user".id = author), "comment", created FROM public."comment" WHERE article = $1`
	rows, err := db.pool.Query(context.Background(), sql, articleId)
	if err != nil {
		return nil, err
	}

	var comments []Comment

	defer rows.Close()
	for rows.Next() {
		var id int
		var root *int
		var author, comment string
		var created time.Time
		err = rows.Scan(&id, &root, &author, &comment, &created)
		if err != nil {
			return nil, err
		}
		comments = append(comments, Comment{
			root,
			id,
			author,
			comment,
			created,
			nil,
		})
	}

	return commentsR(nil, comments), err
}

func (db *PgDB) NewComment(authorId, articleId int, text string, root *int) error {
	panic("implement me")
}

func (db *PgDB) UpdateComment(authorId, commentId int, text string) error {
	panic("implement me")
}

func (db *PgDB) DeleteComment(authorId, commentId int) error {
	panic("implement me")
}
