package api

import (
	"net/http"
	"encoding/json"
	"io"
	"fmt"

	"github.com/curator4/pokedexcli/internal/pokecache"
)

type Pokemon struct {
	Name string
}


type PokemonEncounter struct {
	Pokemon Pokemon
}

type AreaData struct {
	Name string
	Pokemon_encounters []PokemonEncounter
}

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

func GetPage[T any](url string, cache *pokecache.Cache) (*Page[T], error) {

	body, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		cache.Add(url, body)
	}

	var page Page[T]
	err := json.Unmarshal(body, &page)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func GetAreaPokemon(area string, cache *pokecache.Cache) (*AreaData, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"
	url := endpoint + area

	body, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return nil, fmt.Errorf("API returned status %d for URL: %s", res.StatusCode, url)
		}

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if len(body) == 0 {
			return nil, fmt.Errorf("empty response from API for URL: %s", url)
		}

		cache.Add(url, body)
	}

	var areaData AreaData
	err := json.Unmarshal(body, &areaData)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error for %s: %w (response: %s)", url, err, string(body))
	}

	return &areaData, nil
}
