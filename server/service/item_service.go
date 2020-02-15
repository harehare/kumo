package service

import (
	"context"

	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/repository"
)

type ItemService interface {
	List(ctx context.Context, userID string, pageNo, pageSize int) (*[]model.Item, error)
	Get(ctx context.Context, userID, itemID string) (*model.Item, error)
	Delete(ctx context.Context, userID, itemID string) error
	Save(ctx context.Context, item *model.Item) error
}

type ItemServiceImpl struct {
	repository repository.ItemRepository
}

func NewItemService(repository repository.ItemRepository) ItemService {
	return ItemServiceImpl{repository: repository}
}

func (i ItemServiceImpl) List(ctx context.Context, userID string, pageNo, pageSize int) (*[]model.Item, error) {
	offset := (pageNo - 1) * pageSize
	limit := pageNo * pageSize
	return i.repository.FindByUserId(ctx, userID, offset, limit)
}

func (i ItemServiceImpl) Get(ctx context.Context, userID, itemID string) (*model.Item, error) {
	return i.repository.FindOne(ctx, userID, itemID)
}

func (i ItemServiceImpl) Delete(ctx context.Context, userID, itemID string) error {
	return i.repository.Delete(ctx, userID, itemID)
}

func (i ItemServiceImpl) Save(ctx context.Context, item *model.Item) error {
	return i.repository.Save(ctx, item)
}
