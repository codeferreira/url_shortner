package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
	"url_shortner/api"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute", "error", err)
		os.Exit(1)
	}

	slog.Info("system offline")
}

func run() error {
	db := make(map[string]string)
	handler := api.NewHandler(db)

	s := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
