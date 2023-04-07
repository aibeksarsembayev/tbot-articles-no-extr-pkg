package storage

import (
	"context"
	"errors"
	"time"
)

type Storage interface {
	// article api
	Desert(ctx context.Context, a []*ArticleAPIRequestDTO) error

	// author
	DesertAuthor(ctx context.Context, c []*Author) error

	// category
	DesertCategory(ctx context.Context, c []*Category) error

	// article old version
	// Upsert(ctx context.Context, a []*ArticleAPI) error

	GetByCategory(ctx context.Context, categoryID int) ([]*Article, error)
	GetByAuthor(ctx context.Context, userID int) ([]*Article, error)
	GetAllAPI(ctx context.Context) ([]*Article, error)
	GetCategory(ctx context.Context) ([]*Category, error)
	GetAuthor(ctx context.Context) ([]*Author, error)
	GetLatest(ctx context.Context) ([]*Article, error)
	GetLatestMonth(ctx context.Context) ([]*Article, error)
}

var ErrNoArticles = errors.New("no articles")

// ArticleAPI represents article to fetch from API
type ArticleAPI struct {
	ID            int       `json:"article_api_id" db:"id"`
	ArticleID     int       `json:"id" db:"article_id"`
	UserID        int       `json:"user_id" db:"user_id"`
	CategoryID    int       `json:"category_id" db:"category_id"`
	Category      string    `json:"category_name" db:"category"`
	Title         string    `json:"title" db:"title"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	AuthorName    string    `json:"first_name" db:"author_first_name"`
	AuthorSurname string    `json:"last_name" db:"author_last_name"`
	URL           string    `json:"url" db:"url"`
}

// ArticleAuthorDTO represents request artcile author object
type ArticleAuthorDTO struct {
	UserID        int    `json:"user_id" db:"user_id"`
	AuthorName    string `json:"first_name" db:"author_first_name"`
	AuthorSurname string `json:"last_name" db:"author_last_name"`
}

// ArticleCategoryDTO represents request DTO article category DTO
type ArticleCategoryDTO struct {
	Category string `db:"category"`
}

// ArticleV2
// Article represents article object
type Article struct {
	ArticleID int       `json:"id" db:"article_id"`
	Title     string    `json:"title" db:"title"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Article represents article to fetch from API
type ArticleAPIRequestDTO struct {
	ArticleID int         `json:"id" db:"article_id"`
	Category  []*Category `json:"category_name" db:"category"`
	Title     string      `json:"title" db:"title"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	Author    []*Author   `json:"authors" db:"authors"`
	URL       string      `json:"url" db:"url"`
}

// Author represents author object
type Author struct {
	UserID        int    `json:"id" db:"user_id"`
	AuthorName    string `json:"first_name" db:"author_first_name"`
	AuthorSurname string `json:"last_name" db:"author_last_name"`
}

// Category represents category object
type Category struct {
	ID       int    `json:"id" db:"category_id"`
	Category string `json:"name" db:"category_name"`
}
