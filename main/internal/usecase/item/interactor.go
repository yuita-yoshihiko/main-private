package item

import (
	"context"

	"github.com/google/uuid"

	domain "main-private/main/internal/domain/item"
)

type Interactor struct {
	repo domain.Repository
}

func NewInteractor(repo domain.Repository) *Interactor {
	return &Interactor{repo: repo}
}

func (i *Interactor) CreateItem(ctx context.Context, name, description string) (domain.Item, error) {
	newItem, err := domain.NewItem(name, description)
	if err != nil {
		return domain.Item{}, err
	}
	return i.repo.Create(ctx, newItem)
}

func (i *Interactor) GetItem(ctx context.Context, id uuid.UUID) (domain.Item, error) {
	return i.repo.FindByID(ctx, id)
}

func (i *Interactor) ListItems(ctx context.Context) ([]domain.Item, error) {
	return i.repo.FindAll(ctx)
}

func (i *Interactor) UpdateItem(ctx context.Context, id uuid.UUID, name, description string) (domain.Item, error) {
	existing, err := i.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Item{}, err
	}

	existing.Name = name
	existing.Description = description

	if err := existing.Validate(); err != nil {
		return domain.Item{}, err
	}

	return i.repo.Update(ctx, existing)
}

func (i *Interactor) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return i.repo.Delete(ctx, id)
}
