package postgres

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/config"
)

var testStorage *Storage

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	// pseudo-code, some implementation excluded:
	//
	// 1. create test.db if it does not exist
	// 2. run our DDL statements to create the required tables if they do not exist
	// 3. run our tests
	// 4. truncate the test db tables

	// load configs
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Println(conf)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s%s/%s", conf.DBTest.DBUser, conf.DBTest.DBPass, conf.DBTest.DBHost, conf.DBTest.DBPort, conf.DBTest.DBName)

	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return -1, fmt.Errorf("postgres: %w", err)
	}

	testStorage = new(db)

	// truncate all test data after the tests are run
	defer func() {
		for _, t := range []string{"article", "user", "category", "article_user", "article_category"} {
			query := fmt.Sprintf("DELETE FROM %s", t)
			_, _ = db.Exec(query)
		}
		db.Close()
	}()

	return m.Run(), nil
}

func new(db *sqlx.DB) *Storage {
	return &Storage{
		dbpool: db,
	}
}
