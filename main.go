package main

import (
	"context"
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
	defer log.L.Sync()

	_ = tl.InitClient(ctx, *cfg)

	srv := server.NewServer(cfg)
	if err := srv.Run(cfg.Server.GetAddr()); err != nil {
		log.L.Fatal("Failed to run server", zap.Error(err))
	}
}
