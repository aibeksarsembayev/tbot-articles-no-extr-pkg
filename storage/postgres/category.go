package postgres

import (
	"context"

	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
)

// CreateCategory ...
func (s *Storage) DesertCategory(ctx context.Context, c []*storage.Category) error {
	_, err := s.dbpool.Exec(`DELETE FROM "category"`)
	if err != nil {
		return err
	}
	_, err = s.dbpool.NamedExec(`INSERT INTO "category" (category_id, category_name)
	VALUES (:category_id, :category_name)`, c)
	if err != nil {
		return err
	}
	return nil
}

// GetCategoryAPI ...
func (s *Storage) GetCategory(ctx context.Context) ([]*storage.Category, error) {
	c := []*storage.Category{}
	err := s.dbpool.Select(&c, `SELECT category_id, category_name FROM "category"`)
	if err != nil {
		return []*storage.Category{}, err
	}
	return c, nil
}
