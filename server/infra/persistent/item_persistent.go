package persistent

import (
	"context"
	"time"

	"github.com/harehare/kumo/repository"

	"cloud.google.com/go/firestore"
	"github.com/harehare/kumo/model"
	"google.golang.org/api/iterator"
)

type ItemPersistence struct {
	client *firestore.Client
}

func NewItemPersistence(client *firestore.Client) repository.ItemRepository {
	return &ItemPersistence{client: client}
}

func mapToItem(ctx context.Context, data map[string]interface{}) (*model.Item, error) {
	userID := data["UserID"].(string)
	name := data["Name"].(string)
	description := data["Description"].(string)
	URL := data["URL"].(string)
	selector := data["Selector"].(string)

	notifyDocRef, err := data["Notification"].(*firestore.DocumentRef).Get(ctx)

	if err != nil {
		return nil, err
	}

	notification := mapToNotification(notifyDocRef.Data())

	timestamp := data["Time"].(time.Time)
	return &model.Item{
		UserID:       userID,
		Name:         name,
		Description:  description,
		URL:          URL,
		Selector:     selector,
		Notification: notification,
		Timestamp:    timestamp,
	}, nil
}

func (i *ItemPersistence) FindOne(ctx context.Context, userID, itemID string) (*model.Item, error) {
	fields, err := i.client.Collection("users").Doc(userID).Collection("items").Doc(itemID).Get(ctx)

	if err != nil {
		return nil, err
	}

	data := fields.Data()
	item, err := mapToItem(ctx, data)

	if err != nil {
		return nil, err
	}

	item.UserID = userID

	return item, err
}

func (i *ItemPersistence) FindByUserId(ctx context.Context, userID string, offset, limit int) (*[]model.Item, error) {
	var items []model.Item
	iter := i.client.Collection("users").Doc(userID).Collection("items").Offset(offset).Limit(limit).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		data := doc.Data()

		item, err := mapToItem(ctx, data)

		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	return &items, nil
}

func (i *ItemPersistence) Delete(ctx context.Context, userID, itemID string) error {
	_, err := i.client.Collection("users").Doc(userID).Collection("items").Doc(itemID).Delete(ctx)

	return err
}

func (i *ItemPersistence) Save(ctx context.Context, item *model.Item) error {

	resultID := model.NewResultID(item.URL, item.Selector)
	item.Notification.PageURL = item.URL
	item.Notification.Selector = item.Selector
	_, err := i.client.Collection("notify").Doc(resultID.String()).Collection("dest").Doc(item.UserID).Set(ctx, notificationToMap(item.Notification))

	if err != nil {
		return err
	}

	notifyRef := i.client.Collection("notify").Doc(resultID.String()).Collection("dest").Doc(item.UserID)

	_, err = i.client.Collection("users").Doc(item.UserID).Collection("items").Doc(resultID.String()).Set(ctx, map[string]interface{}{
		"UserID":       item.UserID,
		"Name":         item.Name,
		"Description":  item.Description,
		"URL":          item.URL,
		"Selector":     item.Selector,
		"Notification": notifyRef,
		"Time":         firestore.ServerTimestamp,
	})
	if err != nil {
		return err
	}

	return err
}
