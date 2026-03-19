package stats

import (
	"time"

	"github.com/google/uuid"
)

type Stats struct {
	ID         uuid.UUID
	TotalItems int
	CreatedAt  time.Time
}
