package persistent

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/repository"
	"google.golang.org/api/iterator"
)

type ResultPersistence struct {
	client *firestore.Client
}

func NewResultPersistence(client *firestore.Client) repository.ResultRepository {
	return &ResultPersistence{client: client}
}

func mapToResult(data map[string]interface{}) *model.Result {
	URL := data["URL"].(string)
	selector := data["Selector"].(string)
	text := data["Text"].(string)
	err := errors.New(data["Err"].(string))
	timestamp := data["Time"].(time.Time)
	return &model.Result{
		ID:        model.NewResultID(URL, selector),
		URL:       &URL,
		Selector:  &selector,
		Text:      &text,
		Err:       &err,
		Timestamp: timestamp,
	}
}

func (r *ResultPersistence) FindAll(ctx context.Context) (*[]model.Result, error) {
	var results []model.Result
	iter := r.client.Collection("result").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		results = append(results, *mapToResult(data))
	}
	return &results, nil
}

func (r *ResultPersistence) FindOne(ctx context.Context, URL, selector string) (*model.Result, error) {
	resultID := model.NewResultID(URL, selector)

	fields, err := r.client.Collection("result").Doc(resultID.String()).Get(ctx)

	if err != nil {
		return nil, err
	}

	data := fields.Data()
	resultURL := data["URL"].(string)
	resultSelector := data["Selector"].(string)
	text := data["Text"].(string)
	var resultError error
	if data["Err"] != nil {
		resultError = errors.New(data["Err"].(string))
	}
	timestamp := data["Time"].(time.Time)

	return &model.Result{
		URL:       &resultURL,
		Selector:  &resultSelector,
		Text:      &text,
		Err:       &resultError,
		Timestamp: timestamp,
	}, nil
}

func (r *ResultPersistence) Save(ctx context.Context, result *model.Result) error {
	var errString string

	if result.Err != nil {
		err := *result.Err
		errString = err.Error()
	}

	resultID := model.NewResultID(*result.URL, *result.Selector)

	_, err := r.client.Collection("result").Doc(resultID.String()).Set(ctx, map[string]interface{}{
		"URL":      result.URL,
		"Selector": result.Selector,
		"Text":     result.Text,
		"Err":      errString,
		"Time":     firestore.ServerTimestamp,
	})
	if err != nil {
		return err
	}

	return nil
}
