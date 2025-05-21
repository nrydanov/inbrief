package internal

import (
	"github.com/zelenin/go-tdlib/client"
	"github.com/redis/go-redis/v9"
)

type AppState struct {
	TlClient    *client.Client
	RedisClient *redis.Client
}
