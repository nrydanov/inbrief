package config

import (
	"context"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type TelegramConfig struct {
	ApiHash string `env:"API_HASH"`
	ApiId   string `env:"API_ID"`
	Session string `env:"SESSION"`
}

type ServerConfig struct {
	Host string `env:"HOST, default=localhost"`
	Port string `env:"PORT, default=8080"`
}

type Config struct {
	Debug    bool           `env:"DEBUG, default=true"`
	Telegram TelegramConfig `env:", prefix=TELEGRAM_"`
	Server   ServerConfig   `env:", prefix=SERVER_"`
}

func Load(ctx context.Context) (*Config, error) {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Debug {
		log.Printf("Loaded config: %#v", cfg)
	}

	return &cfg, nil
}
