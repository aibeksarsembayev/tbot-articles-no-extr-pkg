package postgres

import (
	"context"

	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
)

// CreateAuthor ...
func (s *Storage) DesertAuthor(ctx context.Context, a []*storage.Author) error {
	_, err := s.dbpool.Exec(`DELETE FROM "user"`)
	if err != nil {
		return err
	}
	_, err = s.dbpool.NamedExec(`INSERT INTO "user" (user_id, author_first_name, author_last_name)
	VALUES (:user_id, :author_first_name, :author_last_name)`, a)
	if err != nil {
		return err
	}
	return nil
}

// GetAuthorAPI ...
func (s *Storage) GetAuthor(ctx context.Context) ([]*storage.Author, error) {
	a := []*storage.Author{}
	err := s.dbpool.Select(&a, `SELECT user_id, author_first_name, author_last_name FROM "user"`)
	if err != nil {
		return []*storage.Author{}, err
	}
	return a, nil
}
