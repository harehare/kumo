package model

import (
	"strings"
	 "encoding/base64"
)

type ResultID struct {
	id string
}

func NewResultID(URL, selector string) *ResultID {
	return &ResultID{id: base64.StdEncoding.EncodeToString([]byte(URL + "," + selector))}
}

func (r *ResultID) FromString(id string) *ResultID {
	return &ResultID{id: id}
}

func (r *ResultID) String() string {
	return r.id
}

func (r *ResultID) Extract() (*string, *string) {
	b, err := base64.StdEncoding.DecodeString(r.id)

	if err != nil {
		return nil, nil
	}

	tokens := strings.Split(",", string(b))

	return &tokens[0], &tokens[1]
}
