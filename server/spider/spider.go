package spider

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/temoto/robotstxt"

	"github.com/harehare/kumo"
	"github.com/harehare/kumo/logger"
	"github.com/harehare/kumo/model"
)

type Spider struct {
	userAgent    string
	allowDomain  string
	delay        time.Duration
	scraperQueue chan Scraper
	robotsTxt    *robotstxt.RobotsData
	parallelism  int
}

func NewSpider(userAgent, allowDomain string, delay time.Duration, parallelism int) *Spider {
	return &Spider{
		userAgent:    userAgent,
		allowDomain:  allowDomain,
		delay:        delay,
		scraperQueue: nil,
		robotsTxt:    nil,
		parallelism:  parallelism,
	}
}

func GetURL(URL string) (*url.URL, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Spider) checkDomain(url string) error {
	u, err := GetURL(url)
	if err != nil {
		return err
	}

	if u.Hostname() != s.allowDomain {
		return kumo.ErrNotAllowDomain
	}

	return nil
}

func (s *Spider) checkRobots(u *url.URL) error {
	logger.Logger.Debug("Check robotx.txt: " + u.Hostname())
	var robot *robotstxt.RobotsData
	if s.robotsTxt == nil {
		resp, err := http.Get(u.Scheme + "://" + u.Host + "/robots.txt")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		robot, err = robotstxt.FromResponse(resp)
		if err != nil {
			return err
		}
		s.robotsTxt = robot
	}

	uaGroup := s.robotsTxt.FindGroup(s.userAgent)
	if uaGroup == nil {
		return nil
	}

	if !uaGroup.Test(u.EscapedPath()) {
		return kumo.ErrRobotsTxtBlocked
	}
	logger.Logger.Debug("Allow robotx.txt")
	return nil
}

func (s *Spider) Start(wg *sync.WaitGroup) <-chan *model.Result {
	logger.Logger.Debug("Slart crawling for " + s.allowDomain)
	s.scraperQueue = make(chan Scraper, s.parallelism)
	c := make(chan *model.Result)
	go func() {
		for scraper := range s.scraperQueue {
			logger.Logger.Debug("Start scraping for " + scraper.URL)
			err := s.checkDomain(scraper.URL)

			if err != nil {
				c <- model.NewFailureResult(scraper.URL, scraper.Selector, err)
				wg.Done()
				continue
			}

			url, err := GetURL(scraper.URL)

			if err != nil {
				c <- model.NewFailureResult(scraper.URL, scraper.Selector, err)
				wg.Done()
				continue
			}

			err = s.checkRobots(url)

			if err != nil {
				c <- model.NewFailureResult(scraper.URL, scraper.Selector, err)
				wg.Done()
				continue
			}

			result, err := scraper.Do()

			if err != nil {
				c <- model.NewFailureResult(scraper.URL, scraper.Selector, err)
				wg.Done()
				time.Sleep(s.delay)
				continue
			}

			c <- model.NewSuccessResult(
				scraper.URL,
				scraper.Selector,
				*result,
			)

			wg.Done()
			logger.Logger.Debug("End scraping")
			time.Sleep(s.delay)
		}

		close(c)
		logger.Logger.Debug("End crawling")
	}()

	return c
}

func (s *Spider) AddScraper(scraper *Scraper) {
	if s.scraperQueue == nil {
		s.scraperQueue = make(chan Scraper, s.parallelism)
	}
	s.scraperQueue <- *scraper
}

func (s *Spider) Exit() {
	close(s.scraperQueue)
	s.scraperQueue = nil
}

type Scraper struct {
	URL      string
	Selector string
}

func NewScraper(url, selector string) *Scraper {
	return &Scraper{URL: url, Selector: selector}
}

func (s *Scraper) Do() (*string, error) {
	res, err := http.Get(s.URL)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	result := doc.Find(s.Selector).Text()

	return &result, err
}
