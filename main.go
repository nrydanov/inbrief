package main

import (
	"context"
	"dsc/inbrief/scraper/internal"
	"dsc/inbrief/scraper/pkg/log"
	defaultlog "log"

	"dsc/inbrief/scraper/config"
	"dsc/inbrief/scraper/internal/server"
	"dsc/inbrief/scraper/internal/tl"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load(ctx)
	if err != nil {
		defaultlog.Fatalf("Failed to load config: %v", err)
	}

	err = log.InitLogger()
	if err != nil {
		defaultlog.Fatalf("Failed to init logger: %v", err)
	}
	defer zap.L().Sync()

	tlClient := tl.InitClient(ctx, *cfg)
	state := internal.AppState{
		TlClient: tlClient,
	}

	server.StartServer(cfg, &state)
}
