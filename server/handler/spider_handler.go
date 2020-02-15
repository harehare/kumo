package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/harehare/kumo"
	"github.com/harehare/kumo/logger"
	"github.com/harehare/kumo/model"
	"github.com/harehare/kumo/spider"
)

// TODO:
const UserAgent = "Kumo spider"

var sn = semaphore.NewWeighted(1)

var spiderPool = sync.Pool{
	New: func() interface{} {
		return &map[string]*spider.Spider{}
	},
}

func merge(cs ...<-chan *model.Result) <-chan *model.Result {
	var wg sync.WaitGroup
	out := make(chan *model.Result)
	output := func(c <-chan *model.Result) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

type SpiderResponse struct {
	Count int `json:"count"`
}

func scraping(scraper *spider.Scraper) (*model.Result, error) {
	text, err := scraper.Do()

	if err != nil {
		return nil, err
	}

	result := model.NewSuccessResult(scraper.URL, scraper.Selector, *text)

	return result, nil
}

func createScraper(param *url.URL) (*spider.Scraper, error) {
	var URL string
	var selector string

	URLQuery, ok := param.Query()["url"]

	if !ok || len(URLQuery[0]) < 1 {
		return nil, kumo.ErrInvalidParameter
	}
	URL = URLQuery[0]

	selectorQuery, ok := param.Query()["selector"]

	if !ok || len(selectorQuery[0]) < 1 {
		return nil, kumo.ErrInvalidParameter
	}

	selector = selectorQuery[0]

	return spider.NewScraper(URL, selector), nil
}

func Entry(w http.ResponseWriter, r *http.Request) *httpError {
	scraper, err := createScraper(r.URL)
	if err != nil {
		return &httpError{Err: err, Msg: "Invalid parameter.", StatusCode: http.StatusBadRequest}
	}

	result, err := scraping(scraper)

	if err != nil {
		return &httpError{Err: nil, Msg: "Error scraping.", StatusCode: http.StatusBadRequest}
	}

	_, err = ResultService.Save(r.Context(), result)

	res, err := json.Marshal(result)
	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	w.Write(res)
	return nil
}

func Test(w http.ResponseWriter, r *http.Request) *httpError {
	scraper, err := createScraper(r.URL)

	if err != nil {
		return &httpError{Err: nil, Msg: "Error scraping.", StatusCode: http.StatusBadRequest}
	}

	result, err := scraping(scraper)

	if err != nil {
		return &httpError{Err: nil, Msg: "Error scraping.", StatusCode: http.StatusBadRequest}
	}

	res, err := json.Marshal(result)
	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	w.Write(res)
	return nil
}

func Crawl(w http.ResponseWriter, r *http.Request) *httpError {
	wg := &sync.WaitGroup{}
	results, err := ResultService.List(r.Context())

	if err != nil {
		return &httpError{Err: err, Msg: "Error fetch result list", StatusCode: http.StatusInternalServerError}
	}

	if !sn.TryAcquire(1) {
		return &httpError{Err: nil, Msg: "Crawl is already running.", StatusCode: http.StatusBadRequest}
	}

	go func() {
		defer sn.Release(1)
		ctx := context.Background()
		logger.Logger.Debug("Start crawling")
		spiders := *spiderPool.Get().(*map[string]*spider.Spider)
		var resultChan []<-chan *model.Result

		go func() {
			wg.Wait()
			for _, spider := range spiders {
				spider.Exit()
			}
		}()

		wg.Add(len(*results) * 2)

		for _, result := range *results {
			url, err := spider.GetURL(*result.URL)
			if err != nil {
				continue
			}
			hostName := url.Hostname()

			if _, ok := spiders[hostName]; !ok {
				spiders[hostName] = spider.NewSpider(UserAgent, hostName, 5*time.Second, 2)
				resultChan = append(resultChan, spiders[hostName].Start(wg))
			}

			spiders[hostName].AddScraper(spider.NewScraper(*result.URL, *result.Selector))
		}

		spiderPool.Put(&spiders)
		updater := spider.NewUpdater(ResultService, NotificationService)

		for result := range merge(resultChan...) {
			err := updater.Update(ctx, result)

			if err != nil {
				// TODO:
			}
			wg.Done()
		}
	}()

	res, err := json.Marshal(&SpiderResponse{Count: len(*results)})
	if err != nil {
		return &httpError{Err: err, Msg: "Error decode json", StatusCode: http.StatusInternalServerError}
	}

	w.Write(res)
	return nil
}
