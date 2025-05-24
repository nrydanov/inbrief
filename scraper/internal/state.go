package internal

import (
	"github.com/aws/aws-sdk-go/service/s3"
	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/redis/go-redis/v9"
	"github.com/zelenin/go-tdlib/client"
)

type ChannelState struct {
	ServerCh   chan *pb.Message
	ListenerCh chan *pb.Message
}

type AppState struct {
	TlClient    *client.Client
	Listener    *client.Listener
	RedisClient *redis.Client
	Channels    *ChannelState
	S3Client    *s3.S3
}

func (s *AppState) Close() {
	s.Listener.Close()
	s.TlClient.Close()
	s.RedisClient.Close()
}
