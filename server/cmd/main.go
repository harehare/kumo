package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"context"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/negroni"

	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"google.golang.org/api/option"

	"github.com/harehare/kumo/handler"
	"github.com/harehare/kumo/infra/persistent"
	"github.com/harehare/kumo/logger"
	"github.com/harehare/kumo/middleware"
	"github.com/harehare/kumo/service"
)

type Env struct {
	Host        string `envconfig:"API_HOST"`
	Port        string `envconfig:"PORT"`
	Credentials string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS_JSON"`
}

func Run() int {
	var env Env
	envconfig.Process("kumo", &env)

	ctx := context.Background()
	b, err := base64.StdEncoding.DecodeString(env.Credentials)

	if err != nil {
		return 1
	}

	opt := option.WithCredentialsJSON(b)
	app, err := firebase.NewApp(ctx, nil, opt)

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return 1
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		log.Fatalln(err)
		return 1
	}

	// TODO: move
	resultRepository := persistent.NewResultPersistence(client)
	notifyRepository := persistent.NewNotifycationPersistence(client)
	itemRepository := persistent.NewItemPersistence(client)
	historyRepository := persistent.NewHistoryPersistence(client)

	resultService := service.NewResultService(resultRepository, historyRepository)
	notifyService := service.NewNotificationService(notifyRepository)
	itemService := service.NewItemService(itemRepository)

	// TODO: DI
	handler.ItemService = itemService
	handler.ResultService = resultService
	handler.NotificationService = notifyService
	logger.Logger = logger.NewLogger()

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"status\": \"OK\"}")
	})

	itemBase := mux.NewRouter()
	r.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AuthMiddleware(app)),
		negroni.Wrap(itemBase),
	))
	itemRoute := itemBase.PathPrefix("/api").Subrouter()
	itemRoute.Methods("GET").Path("/items").HandlerFunc(handler.AppHandler(handler.List).ServeHTTP)
	itemRoute.Methods("GET").Path("/items/{ID}").HandlerFunc(handler.AppHandler(handler.Get).ServeHTTP)
	itemRoute.Methods("DELETE").Path("/items/{ID}").HandlerFunc(handler.AppHandler(handler.Delete).ServeHTTP)
	itemRoute.Methods("POST").Path("/items").HandlerFunc(handler.AppHandler(handler.Save).ServeHTTP)

	spiderBase := mux.NewRouter()
	// TODO: Auth
	r.PathPrefix("/spider").Handler(negroni.New(
		negroni.Wrap(spiderBase)))
	spiderRoute := spiderBase.PathPrefix("/spider").Subrouter()
	spiderRoute.Methods("GET").Path("/entry").HandlerFunc(handler.AppHandler(handler.Entry).ServeHTTP)
	spiderRoute.Methods("GET").Path("/crawl").HandlerFunc(handler.AppHandler(handler.Crawl).ServeHTTP)
	spiderRoute.Methods("GET").Path("/test").HandlerFunc(handler.AppHandler(handler.Test).ServeHTTP)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(middleware.ApiMiddleware))
	n.Use(c)
	n.Use(negronilogrus.NewCustomMiddleware(logrus.InfoLevel, &logrus.JSONFormatter{}, "textusm"))
	n.UseHandler(r)

	s := &http.Server{
		Addr:              fmt.Sprintf(":%s", env.Port),
		Handler:           n,
		ReadTimeout:       8 * time.Second,
		WriteTimeout:      8 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 8 * time.Second,
	}
	err = s.ListenAndServe()

	if err != nil {
		return 1
	}

	return 0
}

func main() {
	os.Exit(Run())
}
