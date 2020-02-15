package repository

import (
	"context"

	"github.com/harehare/kumo/model"
)

type NotificationRepository interface {
	Find(ctx context.Context, URL, selector string) (*[]model.Notification, error)
	FindByUserID(ctx context.Context, URL, selector, UserID string) (*model.Notification, error)
}
