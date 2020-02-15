package service

import (
	"context"
	"net/http"
	"strconv"

	"github.com/harehare/kumo/logger"
	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/repository"
	// "github.com/jordan-wright/email"
)

type NotificationService interface {
	Notify(ctx context.Context, result *model.Result) error
}

type NotificationServiceImpl struct {
	repository repository.NotificationRepository
}

func NewNotificationService(repository repository.NotificationRepository) NotificationService {
	return NotificationServiceImpl{repository: repository}
}

func (n NotificationServiceImpl) Notify(ctx context.Context, result *model.Result) error {
	logger.Logger.Debug("Start notify")
	notifications, err := n.repository.Find(ctx, *result.URL, *result.Selector)

	if err != nil {
		return err
	}

	logger.Logger.Debug("Notification count = " + strconv.Itoa(len(*notifications)))

	for _, notification := range *notifications {
		go webHookNotify(&notification)
	}

	logger.Logger.Debug("End notify")
	return nil
}

// TODO: エラー処理
func webHookNotify(notify *model.Notification) {
	_, err := http.Get(notify.WebHookURL)

	if err != nil {
		// TOOD:
		return
	}

	return
}

func emailNotify(notify *model.Notification) error {
	// if *notify.Email == "" {
	// 	return nil
	// }

	// e := email.NewEmail()
	// e.From = "Jordan Wright <test@gmail.com>"
	// e.To = []string{"test@example.com"}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	// e.Subject = "Awesome Subject"
	// e.Text = []byte("Text Body is, of course, supported!")
	// e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	// e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"))

	return nil
}

func slackNotify(notify *model.Notification) error {
	// TODO:
	return nil
}
