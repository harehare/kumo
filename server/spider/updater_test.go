package spider

import (
	"context"
	"testing"

	"github.com/harehare/kumo/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResultService struct {
	mock.Mock
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockResultService) List(ctx context.Context) (*[]model.Result, error) {
	ret := m.Called(ctx)
	return ret.Get(0).(*[]model.Result), ret.Error(1)
}

func (m *MockResultService) Save(ctx context.Context, result *model.Result) (bool, error) {
	ret := m.Called(ctx, result)
	return ret.Bool(0), ret.Error(1)
}

func (m *MockResultService) AddHistory(ctx context.Context, result *model.Result) error {
	ret := m.Called(ctx, result)
	return ret.Error(0)
}

func (m *MockNotificationService) Notify(ctx context.Context, result *model.Result) error {
	ret := m.Called(ctx, result)
	return ret.Error(0)
}

func TestUpdate(t *testing.T) {
	mockResultService := new(MockResultService)
	mockNotificationService := new(MockNotificationService)

	ctx := context.Background()
	mockChangedResult := model.NewSuccessResult("http://localhost", "changed", "test")
	mockNotChangedResult := model.NewSuccessResult("http://localhost", "not_changed", "test")

	mockResultService.On("AddHistory", ctx, mock.Anything).Return(nil)
	mockResultService.On("Save", ctx,
		mock.MatchedBy(func(r *model.Result) bool {
			return *r.Selector == "changed"
		})).Return(true, nil)
	mockResultService.On("Save", ctx,
		mock.MatchedBy(func(r *model.Result) bool {
			return *r.Selector == "not_changed"
		})).Return(false, nil)

	mockNotificationService.On("Notify", ctx, mock.Anything).Return(nil)

	updater := NewUpdater(mockResultService, mockNotificationService)

	err := updater.Update(ctx, mockChangedResult)
	assert.Nil(t, err)
	mockResultService.AssertCalled(t, "Save", ctx, mockChangedResult)
	mockResultService.AssertCalled(t, "AddHistory", ctx, mockChangedResult)
	mockResultService.AssertNotCalled(t, "List", ctx)
	mockNotificationService.AssertCalled(t, "Notify", ctx, mockChangedResult)

	err = updater.Update(ctx, mockNotChangedResult)
	assert.Nil(t, err)
	mockResultService.AssertCalled(t, "Save", ctx, mockNotChangedResult)
	mockResultService.AssertCalled(t, "AddHistory", ctx, mockNotChangedResult)
	mockResultService.AssertNotCalled(t, "List", ctx)
	mockNotificationService.AssertNotCalled(t, "Notify", ctx, mockNotChangedResult)
}
