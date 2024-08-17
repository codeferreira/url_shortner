package api

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/api/shorten", handlePost(db))
	r.Get("/{code}", handleGet(db))

	return r
}

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  string `json:"data,omitempty"`
}

func sendJSONResponse(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		sendJSONResponse(w, Response{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("failed to write response", "error", err)
		return
	}
}

func handlePost(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSONResponse(w, Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSONResponse(w, Response{Error: "invalid url"}, http.StatusBadRequest)
			return
		}

		code := generateCode()
		db[code] = body.URL

		sendJSONResponse(w, Response{Data: code}, http.StatusCreated)
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() string {
	const n = 8
	byts := make([]byte, n)
	for i := range byts {
		byts[i] = charset[rand.Intn(len(charset))]
	}

	return string(byts)
}

func handleGet(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		url, ok := db[code]
		if !ok {
			http.Error(w, "url not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, url, http.StatusPermanentRedirect)
	}
}
