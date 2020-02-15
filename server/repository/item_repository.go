package repository

import (
	"context"

	"github.com/harehare/kumo/model"
)

type ItemRepository interface {
	FindOne(ctx context.Context, userID, itemID string) (*model.Item, error)
	FindByUserId(ctx context.Context, userID string, offset, limit int) (*[]model.Item, error)
	Delete(ctx context.Context, userID, itemID string) error
	Save(ctx context.Context, item *model.Item) error
}
