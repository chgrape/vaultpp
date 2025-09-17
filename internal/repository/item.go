package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Item struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ItemRepository struct {
	DB *pgxpool.Pool
}

func (r *ItemRepository) GetItems(ctx context.Context) ([]Item, error) {
	var items []Item

	query := `
		SELECT * FROM items
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *ItemRepository) CreateItem(ctx context.Context, item Item) (int, error) {
	var id int

	query := `
		INSERT INTO items(name,description)
		VALUES ($1, $2)
		RETURNING id
	`

	err := r.DB.QueryRow(ctx, query, item.Name, item.Description).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ItemRepository) UpdateItem(ctx context.Context, id int, name string, description string) (int, error) {
	var sets []string
	var args []any
	argPos := 1

	if name != "" {
		sets = append(sets, fmt.Sprintf("name=$%d", argPos))
		args = append(args, name)
		argPos++
	}

	if description != "" {
		sets = append(sets, fmt.Sprintf("description=$%d", argPos))
		args = append(args, description)
		argPos++
	}

	sets = append(sets, "updated_at=NOW()")

	query := fmt.Sprintf(`
		UPDATE items
		SET %s
		WHERE id=$%d
	`, strings.Join(sets, ", "), argPos)
	args = append(args, id)

	tag, err := r.DB.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	if tag.RowsAffected() == 0 {
		return 0, fmt.Errorf("no item with id %d", id)
	}

	return id, nil
}

func (r *ItemRepository) DeleteItem(ctx context.Context, id int) (int, error) {
	query := `
		DELETE FROM items
		WHERE id=$1
	`

	tag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return 0, err
	}

	if tag.RowsAffected() == 0 {
		return 0, fmt.Errorf("no item with id %d", id)
	}
	return id, err
}
