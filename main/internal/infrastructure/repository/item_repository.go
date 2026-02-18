package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	domain "main-private/main/internal/domain/item"
	sqlcdb "main-private/main/internal/infrastructure/db/sqlc"
)

type ItemRepository struct {
	q *sqlcdb.Queries
}

func NewItemRepository(q *sqlcdb.Queries) *ItemRepository {
	return &ItemRepository{q: q}
}

func (r *ItemRepository) Create(ctx context.Context, item domain.Item) (domain.Item, error) {
	row, err := r.q.CreateItem(ctx, sqlcdb.CreateItemParams{
		Name:        item.Name,
		Description: item.Description,
	})
	if err != nil {
		return domain.Item{}, err
	}
	return toDomain(row), nil
}

func (r *ItemRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Item, error) {
	row, err := r.q.GetItem(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Item{}, domain.ErrNotFound
		}
		return domain.Item{}, err
	}
	return toDomain(row), nil
}

func (r *ItemRepository) FindAll(ctx context.Context) ([]domain.Item, error) {
	rows, err := r.q.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]domain.Item, len(rows))
	for i, row := range rows {
		items[i] = toDomain(row)
	}
	return items, nil
}

func (r *ItemRepository) Update(ctx context.Context, item domain.Item) (domain.Item, error) {
	row, err := r.q.UpdateItem(ctx, sqlcdb.UpdateItemParams{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Item{}, domain.ErrNotFound
		}
		return domain.Item{}, err
	}
	return toDomain(row), nil
}

func (r *ItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteItem(ctx, id)
}

func toDomain(row sqlcdb.Item) domain.Item {
	return domain.Item{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
