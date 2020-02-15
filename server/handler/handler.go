package handler

import (
	"net/http"

	"github.com/harehare/kumo/logger"
	"github.com/harehare/kumo/service"
)

type AppHandler func(http.ResponseWriter, *http.Request) *httpError

type httpError struct {
	Err        error
	Msg        string
	StatusCode int
}

var (
	ResultService       service.ResultService
	NotificationService service.NotificationService
	ItemService         service.ItemService
)

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		if err.Err != nil {
			logger.Logger.Error(err.Err.Error())
		}
		http.Error(w, err.Msg, err.StatusCode)
	}
}
