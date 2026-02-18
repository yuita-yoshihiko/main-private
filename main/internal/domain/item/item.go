package item

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("item not found")
	ErrNameRequired = errors.New("item name is required")
)

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewItem(name, description string) (Item, error) {
	i := Item{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := i.Validate(); err != nil {
		return Item{}, err
	}
	return i, nil
}

func (i Item) Validate() error {
	if i.Name == "" {
		return ErrNameRequired
	}
	return nil
}
