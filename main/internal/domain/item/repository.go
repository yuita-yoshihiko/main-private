package item

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, item Item) (Item, error)
	FindByID(ctx context.Context, id uuid.UUID) (Item, error)
	FindAll(ctx context.Context) ([]Item, error)
	Update(ctx context.Context, item Item) (Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
