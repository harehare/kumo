package repository

import (
	"context"

	"github.com/harehare/kumo/model"
)

type ResultRepository interface {
	FindAll(ctx context.Context) (*[]model.Result, error)
	FindOne(ctx context.Context, URL, selector string) (*model.Result, error)
	Save(ctx context.Context, result *model.Result) error
}
