package stats

import "context"

type Repository interface {
	Create(ctx context.Context, totalItems int) (Stats, error)
}
