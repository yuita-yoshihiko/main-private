package stats

import (
	"context"

	domain "main-private/main/internal/domain/stats"
)

type UseCase interface {
	CreateStats(ctx context.Context, totalItems int) (domain.Stats, error)
}
