package item

import (
	"context"

	"github.com/google/uuid"

	domain "main-private/main/internal/domain/item"
)

type UseCase interface {
	CreateItem(ctx context.Context, name, description string) (domain.Item, error)
	GetItem(ctx context.Context, id uuid.UUID) (domain.Item, error)
	ListItems(ctx context.Context) ([]domain.Item, error)
	SearchItems(ctx context.Context, name string) ([]domain.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, name, description string) (domain.Item, error)
	PatchItem(ctx context.Context, id uuid.UUID, name, description *string) (domain.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	GetLatestItem(ctx context.Context) (domain.Item, error)
}
