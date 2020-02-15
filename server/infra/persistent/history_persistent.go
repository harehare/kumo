package persistent

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/repository"
	"google.golang.org/api/iterator"
)

type HistoryPersistence struct {
	client *firestore.Client
}

func NewHistoryPersistence(client *firestore.Client) repository.HistoryRepository {
	return &HistoryPersistence{client: client}
}

func mapToHistory(docID string, data map[string]interface{}) *model.History {
	status := model.NewStatus(data["Status"].(string))
	text := data["Text"].(string)
	timestamp := data["Time"].(time.Time)
	return &model.History{
		ID:        docID,
		Status:    status,
		Text:      text,
		Timestamp: timestamp,
	}
}

func (h *HistoryPersistence) Find(ctx context.Context, itemID string) (*[]model.History, error) {
	var histories []model.History
	iter := h.client.Collection("history").Doc(itemID).Collection("histories").OrderBy("Time", firestore.Desc).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		histories = append(histories, *mapToHistory(doc.Ref.ID, data))
	}
	return &histories, nil
}

func (h *HistoryPersistence) Add(ctx context.Context, itemID string, history *model.History) error {
	_, _, err := h.client.Collection("history").Doc(itemID).Collection("histories").Add(ctx, map[string]interface{}{
		"Status": history.Status.String(),
		"Text":   history.Text,
		"Time":   firestore.ServerTimestamp,
	})
	return err
}

func (h *HistoryPersistence) Rotate(ctx context.Context, itemID string) error {
	histories, err := h.Find(ctx, itemID)

	if err != nil {
		return err
	}

	if len(*histories) > 10 {
		hh := *histories
		deleteHistory := hh[10]
		_, err := h.client.Collection("history").Doc(itemID).Collection("histories").Doc(deleteHistory.ID).Delete(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}
