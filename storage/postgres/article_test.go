package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
	randomizer "github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/tools/random"
)

// func TestStorage_Upsert(t *testing.T) {
// 	args := []*storage.ArticleAPI{}
// 	for i := 0; i < 10; i++ {
// 		arg := randomArticleUpsert()
// 		args = append(args, arg)
// 	}
// 	err := testStorage.Upsert(context.Background(), args)
// 	require.NoError(t, err)
// }

// !!!!!!! need to make random data equal to follow foreign key relation
// func TestStorage_Desert(t *testing.T) {
// 	args := []*storage.ArticleAPIRequestDTO{}
// 	for i := 0; i < 10; i++ {
// 		arg := randomArticle()
// 		args = append(args, arg)
// 	}
// 	err := testStorage.Desert(context.Background(), args)
// 	require.NoError(t, err)
// }

func TestStorage_GetByCategory(t *testing.T) {
	articles := []*storage.Article{}
	category, err := testStorage.GetCategory(context.Background())
	require.NoError(t, err)
	for i := 0; i < len(category); i++ {
		articles, err = testStorage.GetByCategory(context.Background(), category[i].ID)
	}

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetByAuthor(t *testing.T) {
	articles := []*storage.Article{}
	authors, err := testStorage.GetAuthor(context.Background())
	require.NoError(t, err)
	for i := 0; i < len(authors); i++ {
		articles, err = testStorage.GetByAuthor(context.Background(), authors[i].UserID)
	}

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetAllAPI(t *testing.T) {
	articles, err := testStorage.GetAllAPI(context.Background())

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetCategory(t *testing.T) {
	category, err := testStorage.GetCategory(context.Background())
	require.NoError(t, err)

	for _, c := range category {
		require.NotEmpty(t, c)
	}
}

func TestStorage_GetLatest(t *testing.T) {
	articles, err := testStorage.GetLatest(context.Background())

	if errors.Is(err, storage.ErrNoArticles) {
		err = nil
	}

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetLatestMonth(t *testing.T) {
	articles, err := testStorage.GetLatestMonth(context.Background())

	if errors.Is(err, storage.ErrNoArticles) {
		err = nil
	}

	require.NoError(t, err)

	for _, article := range articles {
		require.NotEmpty(t, article)
	}
}

func TestStorage_GetAuthor(t *testing.T) {
	author, err := testStorage.GetAuthor(context.Background())

	require.NoError(t, err)

	if errors.Is(err, storage.ErrNoArticles) {
		err = nil
	}

	for _, a := range author {
		require.NotEmpty(t, a)
	}
}

// func randomArticleUpsert() *storage.ArticleAPI {
// 	return &storage.ArticleAPI{
// 		ArticleID:     int(randomizer.RandomInt(1, 999)),
// 		UserID:        int(randomizer.RandomInt(1, 9)),
// 		CategoryID:    int(randomizer.RandomInt(1, 9)),
// 		Category:      randomizer.RandomString(10),
// 		Title:         randomizer.RandomString(10),
// 		CreatedAt:     randomizer.RandomDate(),
// 		UpdatedAt:     randomizer.RandomDate().Add(10),
// 		AuthorName:    randomizer.RandomString(6),
// 		AuthorSurname: randomizer.RandomString(10),
// 		URL:           randomizer.RandomString(20),
// 	}
// }

func randomArticle() *storage.ArticleAPIRequestDTO {
	return &storage.ArticleAPIRequestDTO{
		ArticleID: int(randomizer.RandomInt(1, 999)),
		Category: []*storage.Category{
			{
				ID:       int(randomizer.RandomInt(1, 9)),
				Category: randomizer.RandomString(10),
			},
		},
		Title:     randomizer.RandomString(10),
		CreatedAt: randomizer.RandomDate(),
		UpdatedAt: randomizer.RandomDate().Add(10),
		Author: []*storage.Author{
			{
				UserID:        int(randomizer.RandomInt(1, 9)),
				AuthorName:    randomizer.RandomString(6),
				AuthorSurname: randomizer.RandomString(10),
			},
		},
		URL: randomizer.RandomString(20),
	}
}
