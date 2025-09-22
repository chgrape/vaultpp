package service

import (
	"context"
	"errors"
	"strings"

	"github.com/chgrape/vaultpp/internal/repository"
)

type ItemService struct {
	Repo *repository.ItemRepository
}

func (s *ItemService) ListItems(ctx context.Context) ([]repository.Item, error) {
	return s.Repo.GetItems(ctx)
}

func (s *ItemService) ListItem(id int, ctx context.Context) (*repository.Item, error) {
	return s.Repo.GetItem(id, ctx)
}

func (s *ItemService) AddItem(item repository.Item, ctx context.Context) (int, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Description) == "" {
		return 0, errors.New("missing required field")
	}
	return s.Repo.CreateItem(ctx, item)
}

func (s *ItemService) EditItem(item repository.Item, id int, ctx context.Context) (int, error) {
	if strings.TrimSpace(item.Name) == "" && strings.TrimSpace(item.Description) == "" {
		return 0, errors.New("at least one field required for update")
	}
	return s.Repo.UpdateItem(ctx, id, item.Name, item.Description)
}

func (s *ItemService) RemoveItem(id int, ctx context.Context) (int, error) {
	return s.Repo.DeleteItem(ctx, id)
}
