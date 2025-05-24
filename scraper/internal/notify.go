package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type Notifier struct {
	ch         <-chan *pb.Message
	s3Client   *s3.S3
	rdb        *redis.Client
	buffer     []*pb.Message
	pubChannel string
	ptr        int

	mu sync.Mutex
}

func NewNotifier(
	ch <-chan *pb.Message,
	s3 *s3.S3,
	rdb *redis.Client,
	pubChannel string,
	bufferSize int,
) *Notifier {
	return &Notifier{
		ch:         ch,
		s3Client:   s3,
		rdb:        rdb,
		pubChannel: pubChannel,
		buffer:     make([]*pb.Message, bufferSize),
	}
}

func (n *Notifier) Listen(ctx context.Context) {
	for msg := range n.ch {
		zap.L().Debug("Received new message", zap.String("text", msg.Text))
		n.mu.Lock()
		n.buffer[n.ptr] = msg
		n.ptr += 1
		timeToNotify := n.ptr == cap(n.buffer)
		n.mu.Unlock()

		if timeToNotify {
			err := n.notify(ctx)
			if err != nil {
				zap.L().Error("Failed to notify", zap.Error(err))
			}
		}
	}
}

func (n *Notifier) NotifyByPeriod(ctx context.Context, period time.Duration) {
	for {
		<-time.After(period)
		err := n.notify(ctx)
		if err != nil {
			zap.L().Error("Failed to notify", zap.Error(err))
		}
	}
}

func (n *Notifier) notify(ctx context.Context) error {
	nMsgs := n.ptr
	if nMsgs == 0 {
		zap.L().Info("Nothing to flush since last time")
		return nil
	}
	zap.L().Info(fmt.Sprintf("Flushing %d messages since last time", nMsgs))

	id := time.Now().UnixNano()

	marshaler := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		Indent:          "  ",
	}

	jsonMessages := make([]json.RawMessage, nMsgs)
	for i, msg := range n.buffer[:nMsgs] {
		jsonData, err := marshaler.Marshal(msg)
		if err != nil {
			zap.L().Error("Failed to marshal proto message", zap.Error(err))
			return err
		}
		jsonMessages[i] = json.RawMessage(jsonData)
	}

	marshalled, err := json.Marshal(jsonMessages)
	if err != nil {
		zap.L().Error("Failed to marshal messages", zap.Error(err))
		return err
	}

	_, err = n.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("inbrief"),
		Key:    aws.String(fmt.Sprintf("%d.json", id)),
		Body:   bytes.NewReader(marshalled),
	})
	if err != nil {
		zap.L().Error("Failed to upload messages to S3", zap.Error(err))
		return err
	}

	zap.L().Info(
		"Successfully flushed messages",
		zap.Int("count", n.ptr),
		zap.String("id", fmt.Sprintf("%d", id)),
	)

	n.rdb.Publish(ctx, n.pubChannel, id)
	n.ptr = 0

	return nil
}
