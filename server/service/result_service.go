package service

import (
	"context"
	"time"

	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/repository"
)

type ResultService interface {
	List(ctx context.Context) (*[]model.Result, error)
	Save(ctx context.Context, result *model.Result) (bool, error)
	AddHistory(ctx context.Context, result *model.Result) error
}

type ResultServiceImpl struct {
	repository        repository.ResultRepository
	historyRepository repository.HistoryRepository
}

func NewResultService(repository repository.ResultRepository, historyRepository repository.HistoryRepository) ResultService {
	return ResultServiceImpl{repository: repository, historyRepository: historyRepository}
}

func (r ResultServiceImpl) List(ctx context.Context) (*[]model.Result, error) {
	return r.repository.FindAll(ctx)
}

func (r ResultServiceImpl) Save(ctx context.Context, result *model.Result) (bool, error) {
	oldResult, _ := r.repository.FindOne(ctx, *result.URL, *result.Selector)
	if oldResult == nil || *result.Text != *oldResult.Text {
		err := r.repository.Save(ctx, result)
		if err != nil {
			return false, err
		}

		return true, nil
	}
	return false, nil
}

func (r ResultServiceImpl) AddHistory(ctx context.Context, result *model.Result) error {
	var addHistory model.History

	if result.Err == nil {
		addHistory = model.History{
			ID:        "",
			Status:    model.Ok,
			Text:      *result.Text,
			Timestamp: time.Now(),
		}
	} else {
		addHistory = model.History{
			ID:        "",
			Status:    model.Err,
			Text:      *result.Text,
			Timestamp: time.Now(),
		}
	}

	err := r.historyRepository.Add(ctx, result.ID.String(), &addHistory)

	if err != nil {
		return err
	}

	err = r.historyRepository.Rotate(ctx, result.ID.String())

	if err != nil {
		return err
	}

	return nil
}
