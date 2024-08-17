package omdb

import (
	"fmt"
	"net/http"
)

type Result struct {
	Search       []SearchResult `json:"search"`
	TotalResults string         `json:"totalResults"`
	Response     string         `json:"response"`
}
type SearchResult struct {
	Title  string `json:"title"`
	Year   string `json:"year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"yype"`
	Poster string `json:"poster"`
}

func Search(query string) (Result, error) {
	resp, err := http.Get("https://www.omdbapi.com/?s=" + query + "&apikey=")
	if err != nil {
		return Result{}, fmt.Errorf("failed to get response: %w", err)
	}

	return Result{}, nil
}
