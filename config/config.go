package config

import (
	"context"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type TelegramConfig struct {
	ApiHash string `env:"API_HASH"`
	ApiId   int32  `env:"API_ID"`
	Session string `env:"SESSION"`
}

type ServerConfig struct {
	Host string `env:"HOST, default=127.0.0.1"`
	Port string `env:"PORT, default=8080"`
}

type RedisConfig struct {
	Host  string `env:"HOST, default=127.0.0.1"`
	Port  string `env:"PORT, default=6379"`
	Topic string `env:"TOPIC, default=inbrief"`
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type Config struct {
	Debug     bool           `env:"DEBUG, default=true"`
	Streaming bool           `env:"STREAMING, default=true"`
	Telegram  TelegramConfig `env:", prefix=TELEGRAM_"`
	Server    ServerConfig   `env:", prefix=SERVER_"`
	Redis     RedisConfig    `env:", prefix=REDIS_"`
}

func (c *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
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
