package spider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/harehare/kumo/logger"
	"github.com/stretchr/testify/assert"
)

var robotsTxtFile = `
User-agent: *
Allow: /allowed
Disallow: /disallowed
`

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`<!DOCTYPE html>
		<html>
			<head>
			<title>test</title>
			</head>
		<body>
			<header></header>
			<h1>h1</h1>
			<main>main</main>
			<footer></footer>
		</body>
		</html>`))
	})

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(robotsTxtFile))
	})

	mux.HandleFunc("/allowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("allowed"))
	})

	mux.HandleFunc("/disallowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("disallowed"))
	})

	return httptest.NewServer(mux)
}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
	logger.Logger = logger.NewLogger()
}

func teardown() {
	// TODO:
}

func TestGetURL(t *testing.T) {
	t.Parallel()
	u, err := GetURL("http://localhost")

	assert.Equal(t, u.Hostname(), "localhost")
	assert.Nil(t, err)
}

func TestCheckDomain(t *testing.T) {
	t.Parallel()
	spider := NewSpider("", "localhost", time.Second, 1)

	err := spider.checkDomain("https://localhost")
	assert.Nil(t, err)

	err = spider.checkDomain("https://localhost1")
	assert.NotNil(t, err)
}

func TestCheckRobots(t *testing.T) {
	t.Parallel()
	spider := NewSpider("", "127.0.0.1", time.Second, 1)

	ts := newTestServer()
	defer ts.Close()

	fmt.Println(ts.URL)

	baseURL, err := GetURL(ts.URL)
	err = spider.checkRobots(baseURL)
	assert.Nil(t, err)

	allowURL, err := GetURL(ts.URL + "/allowed")
	err = spider.checkRobots(allowURL)
	assert.Nil(t, err)

	disallowURL, err := GetURL(ts.URL + "/disallowed")
	err = spider.checkRobots(disallowURL)
	assert.NotNil(t, err)
}

func TestStart(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	wg.Add(5)
	spider := NewSpider("", "127.0.0.1", time.Second, 1)
	ts := newTestServer()
	defer ts.Close()
	defer spider.Exit()

	c := spider.Start(wg)

	spider.AddScraper(NewScraper(ts.URL, "title"))
	titleResult := <-c
	assert.Equal(t, *titleResult.Text, "test")
	assert.Nil(t, titleResult.Err)

	spider.AddScraper(NewScraper(ts.URL, "main"))
	mainResult := <-c
	assert.Equal(t, *mainResult.Text, "main")
	assert.Nil(t, mainResult.Err)

	spider.AddScraper(NewScraper(ts.URL, "header"))
	headerResult := <-c
	assert.Equal(t, *headerResult.Text, "")
	assert.Nil(t, headerResult.Err)

	spider.AddScraper(NewScraper("http://localhost:8080", "h1"))
	domainErrResult := <-c
	assert.NotNil(t, domainErrResult.Err)

	spider.AddScraper(NewScraper(ts.URL+"/disallowed", "title"))
	checkRobotsTxtResult := <-c
	assert.NotNil(t, *checkRobotsTxtResult.Err)
}

func TestScraperDo(t *testing.T) {
	t.Parallel()

	ts := newTestServer()
	defer ts.Close()

	titleScraper := NewScraper(ts.URL, "title")
	title, err := titleScraper.Do()
	assert.Equal(t, *title, "test")
	assert.Nil(t, err)

	headerScraper := NewScraper(ts.URL, "header")
	header, err := headerScraper.Do()
	assert.Equal(t, *header, "")
	assert.Nil(t, err)
}
