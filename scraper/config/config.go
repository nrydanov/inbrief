package config

import (
	"context"
	"fmt"
	"log"
	"time"

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
	Host    string `env:"HOST, default=127.0.0.1"`
	Port    string `env:"PORT, default=6379"`
	Channel string `env:"CHANNEL, default=inbrief"`
}

type S3Config struct {
	Endpoint string `env:"ENDPOINT, default=http://127.0.0.1:9000"`
	Region   string `env:"REGION, default=us-east-1"`
	Username string `env:"USERNAME, default=minioadmin"`
	Password string `env:"PASSWORD, default=minioadmin"`
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type StreamingConfig struct {
	On          bool          `env:"ON, default=true"`
	FlushPeriod time.Duration `env:"FLUSH_PERIOD, default=5s"`
	BatchSize   int           `env:"BATCHSIZE, default=1000"`
}

type Config struct {
	Debug     bool            `env:"DEBUG, default=true"`
	Streaming StreamingConfig `env:", prefix=STREAMING_"`
	Telegram  TelegramConfig  `env:", prefix=TELEGRAM_"`
	Server    ServerConfig    `env:", prefix=SERVER_"`
	Redis     RedisConfig     `env:", prefix=REDIS_"`
	S3        S3Config        `env:", prefix=S3_"`
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
