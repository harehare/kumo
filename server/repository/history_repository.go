package repository

import (
	"context"

	"github.com/harehare/kumo/model"
)

type HistoryRepository interface {
	Find(ctx context.Context, itemID string) (*[]model.History, error)
	Add(ctx context.Context, itemID string, history *model.History) error
	Rotate(ctx context.Context, itemID string) error
}
