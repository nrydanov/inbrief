package tl

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/redis/go-redis/v9"
	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventHandler struct {
	listener   *client.Listener
	rdb        *redis.Client
	s3Client   *s3.S3
	msgBuffer  []*pb.Message
	pubChannel string
	ptr        int
}

func NewEventHandler(
	listener *client.Listener,
	rdb *redis.Client,
	s3 *s3.S3,
	msgBufferSize int,
	publishChannel string,
) *EventHandler {
	return &EventHandler{
		listener:   listener,
		rdb:        rdb,
		s3Client:   s3,
		msgBuffer:  make([]*pb.Message, msgBufferSize),
		pubChannel: publishChannel,
	}
}

func (eh *EventHandler) flush() error {
	zap.L().Info(fmt.Sprintf("Flushing %d messages since last time", eh.ptr))

	uuid := uuid.New()

	eh.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("inbrief_fetcher"),
		Key:    aws.String(uuid.String()),
	})
	eh.ptr = 0

	return nil
}

func (eh *EventHandler) Handle(listener *client.Listener, rdb *redis.Client) {
	for update := range listener.Updates {
		switch msg := update.(type) {
		case *client.UpdateNewMessage:
			eh.newMessageHandler(msg)
		}
	}
}

func (eh *EventHandler) FlushByPeriod(ctx context.Context, period time.Duration) {
	for {
		<-time.After(period)
		eh.flush()
	}
}

func (eh *EventHandler) newMessageHandler(msg *client.UpdateNewMessage) {
	zap.L().Info(
		"New message",
		zap.String("chat_id", fmt.Sprintf(
			"%d",
			msg.Message.ChatId,
		)),
		zap.String("message_id", fmt.Sprintf(
			"%d",
			msg.Message.Id,
		)),
	)
	switch content := msg.Message.Content.(type) {
	case *client.MessageText:
		zap.L().Info(
			"New message text",
			zap.String("text", content.Text.Text),
		)
		eh.msgBuffer[eh.ptr] = &pb.Message{
			Id:     msg.Message.Id,
			Text:   content.Text.Text,
			Ts:     timestamppb.New(time.Unix(int64(msg.Message.Date), 0)),
			ChatId: msg.Message.ChatId,
		}
		eh.ptr += 1

		if eh.ptr == len(eh.msgBuffer) {
			eh.flush()
		}
	}
}
