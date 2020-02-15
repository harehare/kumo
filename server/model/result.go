package model

import (
	"encoding/json"
	"time"
)

type Result struct {
	ID        *ResultID
	URL       *string `json:"url"`
	Selector  *string `json:"selector"`
	Text      *string `json:"text"`
	Err       *error
	Timestamp time.Time
}

func (r Result) MarshalJSON() ([]byte, error) {
	var tmp struct {
		URL      *string `json:"url"`
		Selector *string `json:"selector"`
		Text     *string `json:"text"`
	}
	tmp.URL = r.URL
	tmp.Selector = r.Selector
	tmp.Text = r.Text
	return json.Marshal(&tmp)
}

func NewSuccessResult(url, selector, text string) *Result {
	return &Result{
		ID:        NewResultID(url, selector),
		URL:       &url,
		Selector:  &selector,
		Text:      &text,
		Err:       nil,
		Timestamp: time.Now(),
	}
}

func NewFailureResult(url, selector string, err error) *Result {
	return &Result{
		ID:        NewResultID(url, selector),
		URL:       &url,
		Selector:  &selector,
		Text:      nil,
		Err:       &err,
		Timestamp: time.Now(),
	}
}
