package api

import (
	"net/http"
	"encoding/json"
	"io"
)

type LocationArea struct {
	Id int
	Name string
}

type Page[T any] struct {
	Count int
	Next string
	Previous string
	Results []T
}
func GetPage[T any](url string) (*Page[T], error) {
	
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var page Page[T]
	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}
