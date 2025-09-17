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

func (s *ItemService) AddItem(item repository.Item, ctx context.Context) (int, error) {
	if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Description) == "" {
		return 0, errors.New("Missing required field")
	}
	return s.Repo.CreateItem(ctx, item)
}
