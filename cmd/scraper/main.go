package main

import (
	"context"
	"github.com/nrydanov/inbrief/internal"
	"github.com/nrydanov/inbrief/pkg/log"
	defaultlog "log"

	"github.com/nrydanov/inbrief/config"
	"github.com/nrydanov/inbrief/internal/server"
	"github.com/nrydanov/inbrief/internal/tl"

	"github.com/redis/go-redis/v9"
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	state := internal.AppState{
		TlClient:    tlClient,
		RedisClient: rdb,
	}

	server.StartServer(cfg, &state)
}
