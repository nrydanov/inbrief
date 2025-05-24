package main

import (
	"context"
	defaultlog "log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nrydanov/inbrief/internal"
	"github.com/nrydanov/inbrief/pkg/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nrydanov/inbrief/config"
	"github.com/nrydanov/inbrief/internal/server"
	"github.com/nrydanov/inbrief/internal/tl"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

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
	defer tlClient.Close()

	var rdb *redis.Client
	var s3Client *s3.S3
	if cfg.Streaming.On {
		{
			rdb = redis.NewClient(&redis.Options{
				Addr:     cfg.Redis.GetAddr(),
				Password: "", // no password set
				DB:       0,  // use default DB
			})
			defer rdb.Close()

			if err := rdb.Ping(ctx).Err(); err != nil {
				zap.L().Fatal("Failed to ping Redis", zap.Error(err))
			} else {
				zap.L().Info("Initialized Redis client successfully")
			}
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

	wg := sync.WaitGroup{}

	eventHandler, eventCh := tl.NewEventHandler(
		state.Listener,
		cfg.Streaming.BatchSize,
	)

	notifier := internal.NewNotifier(
		eventCh,
		state.S3Client,
		state.RedisClient,
		cfg.Redis.Channel,
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		notifier.Listen(ctx, cfg.Streaming.BatchSize)
		zap.L().Debug("Notifier is stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		eventHandler.Handle(ctx, state.Listener, state.RedisClient)
		zap.L().Debug("Event handler is stopped")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.StartServer(ctx, cfg, &state)
		zap.L().Debug("RPC server is stopped")
	}()

	wg.Wait()
	zap.L().Info("All workers are stopped, exiting")

}
