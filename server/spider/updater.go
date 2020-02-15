package spider

import (
	"context"

	"github.com/harehare/kumo/logger"
	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/service"
)

type Updater struct {
	resultService       service.ResultService
	notificationService service.NotificationService
}

func NewUpdater(resultService service.ResultService, notificationService service.NotificationService) *Updater {
	return &Updater{resultService: resultService, notificationService: notificationService}
}

func (u *Updater) Update(ctx context.Context, result *model.Result) error {
	logger.Logger.Debug("Start updater")
	changed, err := u.resultService.Save(ctx, result)

	if err != nil {
		result.Err = &err
		u.resultService.AddHistory(ctx, result)
		return err
	}

	if changed {
		logger.Logger.Debug("Changed for " + *result.URL)
		err := u.notificationService.Notify(ctx, result)

		if err != nil {
			result.Err = &err
			u.resultService.AddHistory(ctx, result)
			return err
		}
	}

	u.resultService.AddHistory(ctx, result)
	logger.Logger.Debug("End updater")
	return nil
}
