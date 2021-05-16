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
    id SERIAL,
    login character varying(250),
    pass_hash character(60),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.article
(
    id SERIAL,
    author integer NOT NULL,
    title character varying(250) NOT NULL,
    body text NOT NULL,
    created timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.comment
(
    id SERIAL,
    root integer,
    author integer NOT NULL,
    comment character varying(500) NOT NULL,
    article integer NOT NULL,
    created timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE public.article
    ADD FOREIGN KEY (author)
    REFERENCES public."user" (id) ON DELETE cascade 
    NOT VALID;


ALTER TABLE public.comment
    ADD FOREIGN KEY (author)
    REFERENCES public."user" (id)
    NOT VALID;


ALTER TABLE public.comment
    ADD FOREIGN KEY (article)
    REFERENCES public.article (id) ON DELETE cascade
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
	r, err := db.pool.Exec(context.Background(), `INSERT INTO public."user"(login, pass_hash) VALUES($1, $2)`, login, passHash)
	if err != nil {
		return nil, err
	} else if r.RowsAffected() == 0 {
		return nil, errors.New("zero affected rows")
	}

	var userId int
	err = db.pool.QueryRow(context.Background(), `SELECT id FROM public."user" WHERE login = $1`, login).Scan(&userId)
	if err != nil {
		return nil, err
	}

	return &userId, nil
}

func (db *PgDB) GetUserByLogin(login string) (*User, error) {
	var userId int
	var passHash string
	sql := `SELECT id, pass_hash FROM public."user" WHERE login = $1`
	err := db.pool.QueryRow(context.Background(), sql, login).Scan(&userId, &passHash)
	if err != nil {
		return nil, err
	}

	return &User{userId, login, passHash}, nil
}

func (db *PgDB) GetUserById(userId int) (*User, error) {
	var login, passHash string
	sql := `SELECT login, pass_hash FROM public."user" WHERE id = $1`
	err := db.pool.QueryRow(context.Background(), sql, userId).Scan(&login, &passHash)
	if err != nil {
		return nil, err
	}

	return &User{userId, login, passHash}, nil
}

func (db *PgDB) NewArticle(authorId int, title, text string) (*int, error) {
	timestamp := time.Now()
	sql := `INSERT INTO public.article(author, title, body, created) VALUES($1, $2, $3, $4)`
	r, err := db.pool.Exec(context.Background(), sql, authorId, title, text, timestamp)
	if err != nil {
		return nil, err
	} else if r.RowsAffected() != 1 {
		return nil, errors.New("no rows effected")
	}

	var articleId int
	sql = `SELECT id FROM public.article WHERE author = $1 AND created = $2`
	err = db.pool.QueryRow(context.Background(), sql, authorId, timestamp).Scan(&articleId)
	if err != nil {
		return nil, err
	}

	return &articleId, nil
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

func (db *PgDB) checkAuthor(userId, articleId int) (bool, error) {
	var authorId *int
	sql := `SELECT author FROM public.article WHERE id = $1`
	err := db.pool.QueryRow(context.Background(), sql, articleId).Scan(&authorId)
	if err != nil {
		return false, err
	} else if authorId == nil {
		return false, nil
	} else if *authorId != userId {
		return false, nil
	} else {
		return true, nil
	}
}

func (db *PgDB) UpdateArticle(userId, articleId int, title, body string) error {
	isAuthor, err := db.checkAuthor(userId, articleId)
	if !isAuthor {
		return errors.New("only author can update article")
	}

	sql := `UPDATE public.article SET title = $2, body = $3 WHERE id = $1`
	_, err = db.pool.Exec(context.Background(), sql, articleId, title, body)
	return err
}

func (db *PgDB) DeleteArticle(userId, articleId int) error {
	isAuthor, err := db.checkAuthor(userId, articleId)
	if !isAuthor {
		return errors.New("only author can update article")
	}

	sql := `DELETE FROM public.article WHERE id = $1`
	_, err = db.pool.Exec(context.Background(), sql, articleId)
	return err
}

func commentsR(root *int, comments []Comment) []Comment {
	var parsedComments []Comment

	for _, comment := range comments {
		add := false
		if comment.Root == root {
			add = true
		} else if root != nil && comment.Root != nil {
			if *comment.Root == *root {
				add = true
			}
		}
		if add {
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
	sql := `INSERT INTO public.comment(author, article, comment, root, created) VALUES($1, $2, $3, $4, $5)`
	_, err := db.pool.Query(context.Background(), sql, authorId, articleId, text, root, time.Now())
	return err
}

func (db *PgDB) UpdateComment(authorId, commentId int, text string) error {
	panic("implement me")
}

func (db *PgDB) DeleteComment(authorId, commentId int) error {
	panic("implement me")
}
