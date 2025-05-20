package main

import (
	"context"
	"log"

	"dsc/inbrief/scraper/config"
	"dsc/inbrief/scraper/internal/server"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := server.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
