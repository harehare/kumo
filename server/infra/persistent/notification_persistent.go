package persistent

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/repository"
	"google.golang.org/api/iterator"
)

type NotifycationPersistence struct {
	client *firestore.Client
}

func NewNotifycationPersistence(client *firestore.Client) repository.NotificationRepository {
	return &NotifycationPersistence{client: client}
}

func mapToNotification(data map[string]interface{}) *model.Notification {
	pageURL := data["PageURL"].(string)
	notifySelector := data["Selector"].(string)
	text := data["Text"].(string)
	webHookURL := data["WebHookURL"].(string)
	timestamp := data["Time"].(time.Time)
	return &model.Notification{
		PageURL:    pageURL,
		Selector:   notifySelector,
		Text:       text,
		WebHookURL: webHookURL,
		Timestamp:  timestamp,
	}
}

func notificationToMap(data *model.Notification) *map[string]interface{} {
	return &map[string]interface{}{
		"PageURL":    data.PageURL,
		"Selector":   data.Selector,
		"Text":       data.Text,
		"WebHookURL": data.WebHookURL,
		"Time":       firestore.ServerTimestamp,
	}
}

func (n *NotifycationPersistence) Find(ctx context.Context, URL, selector string) (*[]model.Notification, error) {
	var notifications []model.Notification
	resultID := model.NewResultID(URL, selector)
	iter := n.client.Collection("notify").Doc(resultID.String()).Collection("dest").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()
		notifications = append(notifications, *mapToNotification(data))
	}
	return &notifications, nil
}

func (n *NotifycationPersistence) FindByUserID(ctx context.Context, URL, selector, UserID string) (*model.Notification, error) {
	resultID := model.NewResultID(URL, selector)
	fields, err := n.client.Collection("notify").Doc(resultID.String()).Collection("dest").Doc(UserID).Get(ctx)

	if err != nil {
		return nil, err
	}

	data := fields.Data()
	return mapToNotification(data), nil
}
