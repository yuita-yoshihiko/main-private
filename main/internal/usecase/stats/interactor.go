package stats

import (
	"context"

	domain "main-private/main/internal/domain/stats"
)

type Interactor struct {
	repo domain.Repository
}

func NewInteractor(repo domain.Repository) *Interactor {
	return &Interactor{repo: repo}
}

func (i *Interactor) CreateStats(ctx context.Context, totalItems int) (domain.Stats, error) {
	return i.repo.Create(ctx, totalItems)
}
