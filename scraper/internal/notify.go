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
	pubChannel string
}

func NewNotifier(
	ch <-chan *pb.Message,
	s3 *s3.S3,
	rdb *redis.Client,
	pubChannel string,
) *Notifier {
	return &Notifier{
		ch:         ch,
		s3Client:   s3,
		rdb:        rdb,
		pubChannel: pubChannel,
	}
}

func (n *Notifier) Listen(ctx context.Context, bufferSize int) {

	ticker := time.NewTicker(time.Second * 5)

	buffer := make([]*pb.Message, bufferSize)
	ptr := 0
	flushCh := make(chan []*pb.Message)
	sendSafe := func() {
		newBuffer := make([]*pb.Message, ptr)
		copy(newBuffer, buffer[:ptr])
		flushCh <- newBuffer
		ptr = 0
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	defer wg.Wait()
	defer ticker.Stop()
	defer close(flushCh)
	defer sendSafe()

	go func() {
		defer wg.Done()
		for msgs := range flushCh {
			err := n.notify(msgs)
			if err != nil {
				zap.L().Error("Failed to notify", zap.Error(err))
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendSafe()
		case msg, ok := <-n.ch:
			if !ok {
				return
			}
			zap.L().Debug("Received new message", zap.String("text", msg.Text))
			buffer[ptr] = msg
			ptr += 1

			if ptr == cap(buffer) {
				sendSafe()
			}
		}
	}
}

func (n *Notifier) notify(msgs []*pb.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nMsgs := len(msgs)
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
	for i, msg := range msgs {
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
		zap.Int("count", nMsgs),
		zap.String("id", fmt.Sprintf("%d", id)),
	)

	return n.rdb.Publish(ctx, n.pubChannel, id).Err()

}
