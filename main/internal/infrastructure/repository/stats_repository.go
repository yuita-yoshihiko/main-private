package repository

import (
	"context"

	domain "main-private/main/internal/domain/stats"
	sqlcdb "main-private/main/internal/infrastructure/db/sqlc"
)

type StatsRepository struct {
	q *sqlcdb.Queries
}

func NewStatsRepository(q *sqlcdb.Queries) *StatsRepository {
	return &StatsRepository{q: q}
}

func (r *StatsRepository) Create(ctx context.Context, totalItems int) (domain.Stats, error) {
	row, err := r.q.CreateStats(ctx, int32(totalItems))
	if err != nil {
		return domain.Stats{}, err
	}
	return domain.Stats{
		ID:         row.ID,
		TotalItems: int(row.TotalItems),
		CreatedAt:  row.CreatedAt.Time,
	}, nil
}
