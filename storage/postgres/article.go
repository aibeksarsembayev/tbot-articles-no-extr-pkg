package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
)

type Storage struct {
	dbpool *sqlx.DB
}

func NewDBArticleRepo(dbpool *sqlx.DB) *Storage {
	return &Storage{
		dbpool: dbpool,
	}
}

// Desert API articles delete old data and insert ...
func (s *Storage) Desert(ctx context.Context, a []*storage.ArticleAPIRequestDTO) error {
	// delete all rows with old data
	_, err := s.dbpool.Exec(`DELETE FROM "article";
	DELETE FROM "article_user";
	DELETE FROM "article_category"`)
	if err != nil {
		return err
	}
	// insert new articles
	_, err = s.dbpool.NamedExec(`INSERT INTO "article" (article_id, title, created_at, updated_at, url)
	VALUES (:article_id, :title, :created_at, :updated_at,:url)`, a)
	if err != nil {
		return err
	}

	// insert article vs user and category info
	for _, article := range a {
		for _, auth := range article.Author {
			_, err = s.dbpool.Exec(`INSERT INTO "article_user" (article_id, user_id)
			VALUES ($1, $2)`, article.ArticleID, auth.UserID)
			if err != nil {
				return err
			}
		}
		for _, cat := range article.Category {
			_, err = s.dbpool.Exec(`INSERT INTO "article_category" (article_id, category_id)
			VALUES ($1, $2)`, article.ArticleID, cat.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Upsert API articles ...
// func (s *Storage) Upsert(ctx context.Context, a []*storage.ArticleAPI) error {
// 	_, err := s.dbpool.NamedExec(`INSERT INTO "article_api" (article_id, user_id, category_id, category, title, created_at, updated_at, author_first_name, author_last_name, url )
// 	VALUES (:article_id, :user_id, :category_id, :category, :title, :created_at, :updated_at, :author_first_name, :author_last_name, :url)
// 	ON CONFLICT(article_id)
// 	DO UPDATE SET (user_id, category_id, category, title, created_at, updated_at, author_first_name, author_last_name, url) =
// 	(excluded.user_id, excluded.category_id, excluded.category, excluded.title, excluded.created_at, excluded.updated_at, excluded.author_first_name, excluded.author_last_name, excluded.url)`, a)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetByCategoryAPI ...
func (s *Storage) GetByCategory(ctx context.Context, categoryID int) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT article.article_id, title, created_at, updated_at, url FROM "article"
	INNER JOIN "article_category" ON article_category.article_id = article.article_id
	WHERE article_category.category_id = $1
	ORDER by created_at DESC`, categoryID)
	if err != nil {
		return []*storage.Article{}, err
	}
	return a, nil
}

// GetByAuthorAPI ...
func (s *Storage) GetByAuthor(ctx context.Context, userID int) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT article.article_id, title, created_at, updated_at, url FROM "article"
	INNER JOIN "article_user" ON article_user.article_id = article.article_id
	WHERE article_user.user_id = $1
	ORDER by created_at DESC`, userID)
	if err != nil {
		return []*storage.Article{}, err
	}
	return a, nil
}

// GetAllAPI ...
func (s *Storage) GetAllAPI(ctx context.Context) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article" ORDER by created_at DESC`)
	if err != nil {
		return []*storage.Article{}, err
	}
	return a, nil
}

// GetLatest article for 7 days ...
func (s *Storage) GetLatest(ctx context.Context) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article" 
	WHERE created_at > current_date - interval '7' day
	ORDER by created_at DESC`)
	if err != nil {
		return []*storage.Article{}, err
	}
	if len(a) == 0 {
		return []*storage.Article{}, storage.ErrNoArticles
	}
	return a, nil
}

// GetLatestMonth article for 30 days ...
func (s *Storage) GetLatestMonth(ctx context.Context) ([]*storage.Article, error) {
	a := []*storage.Article{}
	err := s.dbpool.Select(&a, `SELECT * FROM "article" 
	WHERE created_at > current_date - interval '30' day
	ORDER by created_at DESC`)
	if err != nil {
		return []*storage.Article{}, err
	}
	if len(a) == 0 {
		return []*storage.Article{}, storage.ErrNoArticles
	}
	return a, nil
}
