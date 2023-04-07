package apifetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage/postgres"
	"go.uber.org/zap"
)

const (
	apiURL     = "https://sber-invest.kz/article/general/getall"
	articleURL = "https://sber-invest.kz/article/"
)

type Fetcher struct {
	lg      *zap.Logger
	storage *postgres.Storage
	period  time.Duration
}

func New(logger *zap.Logger, storage *postgres.Storage, period time.Duration) *Fetcher {
	// initial fetch
	return &Fetcher{
		lg:      logger,
		storage: storage,
		period:  period,
	}
}

func (f *Fetcher) Fetch() {
	// initial fetch
	err := f.fetching()
	if err != nil {
		f.lg.Sugar().Error(err)
	}
	t := time.NewTicker(f.period * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err = f.fetching()
			if err != nil {
				f.lg.Sugar().Error(err)
			}
		}
	}
}

func (f *Fetcher) fetching() error {
	url := apiURL

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("fetch api: can't wrap request: %w", err)
	}

	// req.Header.Set("User-Agent","test")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch api: can't send request: %w", err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("fetch api: can't read response: %w", err)
	}

	articles := []*storage.ArticleAPIRequestDTO{}
	err = json.Unmarshal(body, &articles)
	if err != nil {
		return fmt.Errorf("fetch api: can't unmarshal response: %w", err)
	}

	// create url based on article_id
	for _, article := range articles {
		article.URL = fmt.Sprintf("%v%v", articleURL, article.ArticleID)
	}

	// for _, a := range articles {
	// 	fmt.Println(a)
	// 	for _, c := range a.Category {
	// 		fmt.Println(c)
	// 	}
	// 	for _, ax := range a.Author {
	// 		fmt.Println(ax)
	// 	}
	// }

	// retrieve author and category data
	categories := []*storage.Category{}
	authors := []*storage.Author{}
	for _, a := range articles {
		for _, c := range a.Category {
			if isContainCat(c, categories) {
				categories = append(categories, c)
			}
		}
		for _, ax := range a.Author {
			if isContainAuth(ax, authors) {
				authors = append(authors, ax)
			}
		}
	}

	// update db with updated author and category data
	err = f.storage.DesertCategory(context.Background(), categories)
	if err != nil {
		return fmt.Errorf("fetch api: can't delete old categories or insert fetched categories into db: %w", err)
	}
	err = f.storage.DesertAuthor(context.Background(), authors)
	if err != nil {
		return fmt.Errorf("fetch api: can't delete old authors or insert fetched authors into db: %w", err)
	}

	// delete old data and insert new articles
	err = f.storage.Desert(context.Background(), articles)
	if err != nil {
		return fmt.Errorf("fetch api: can't delete old data or insert fetched data into db: %w", err)
	}

	// articles2 := []*storage.ArticleAPI{}
	// err = json.Unmarshal(body, &articles2)
	// if err != nil {
	// 	return fmt.Errorf("fetch api: can't unmarshal response: %w", err)
	// }

	// err = f.storage.Upsert(context.Background(), articles2)
	// if err != nil {
	// 	// log.Fatal("fetch api: can't insert or update fetched data into db: ", err)
	// 	return fmt.Errorf("fetch api: can't insert or update fetched data into db: %w", err)
	// }

	f.lg.Info("API fetching was done successfully")
	return nil
}

// isContainCat checks if category already in array structure
func isContainCat(c *storage.Category, categories []*storage.Category) bool {
	for _, category := range categories {
		if category.ID == c.ID {
			return false
		}
	}
	return true
}

// isContainAuth checks if author already in array structure
func isContainAuth(ax *storage.Author, authors []*storage.Author) bool {
	for _, a := range authors {
		if a.UserID == ax.UserID {
			return false
		}
	}
	return true
}
