package internal

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/redis/go-redis/v9"
	"github.com/zelenin/go-tdlib/client"
)

type AppState struct {
	TlClient    *client.Client
	Listener    *client.Listener
	RedisClient *redis.Client
	S3Client    *s3.S3
}
