package main

import (
	"context"
	defaultlog "log"
	"time"

	"github.com/nrydanov/inbrief/internal"
	"github.com/nrydanov/inbrief/pkg/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nrydanov/inbrief/config"
	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
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
	var rdb *redis.Client
	var s3Client *s3.S3
	if cfg.Streaming.On {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.GetAddr(),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		if err := rdb.Ping(ctx).Err(); err != nil {
			zap.L().Fatal("Failed to ping Redis", zap.Error(err))
		} else {
			zap.L().Info("Initialized Redis client successfully")
		}

		{
			session := session.Must(session.NewSession())
			s3Client = s3.New(
				session,
				aws.NewConfig().
					WithRegion(cfg.S3.Region).
					WithCredentials(credentials.NewStaticCredentials(
						cfg.S3.Username,
						cfg.S3.Password,
						"",
					)).WithEndpoint(cfg.S3.Endpoint),
			)

			if _, err := s3Client.ListBuckets(&s3.ListBucketsInput{}); err != nil {
				zap.L().Fatal("Failed to initialize S3 client", zap.Error(err))
			} else {
				zap.L().Info("Initialized S3 client successfully")
			}
		}

		s3Client.Config.S3ForcePathStyle = aws.Bool(true)
	}

	state := internal.AppState{
		TlClient:    tlClient,
		RedisClient: rdb,
		Listener:    tlClient.GetListener(),
		S3Client:    s3Client,
	}

	eventCh := make(chan *pb.Message, 100)

	eventHandler := tl.NewEventHandler(
		state.Listener,
		eventCh,
	)

	notifier := internal.NewNotifier(
		eventCh,
		state.S3Client,
		state.RedisClient,
		cfg.Redis.Channel,
		cfg.Streaming.BatchSize,
	)

	go notifier.Listen(ctx)
	go notifier.NotifyByPeriod(ctx, time.Second*5)
	go eventHandler.Handle(ctx, state.Listener, state.RedisClient)
	server.StartServer(cfg, &state)

}
