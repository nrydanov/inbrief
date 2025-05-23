package tl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"time"
	"unicode/utf16"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/redis/go-redis/v9"
	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
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

func (eh *EventHandler) flush(ctx context.Context) error {
	if eh.ptr == 0 {
		zap.L().Info("Nothing to flush since last time")
		return nil
	}
	zap.L().Info(fmt.Sprintf("Flushing %d messages since last time", eh.ptr))

	id := time.Now().UnixNano()

	marshaler := protojson.MarshalOptions{
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
		Indent:          "  ",
	}

	jsonMessages := make([]json.RawMessage, eh.ptr)
	for i, msg := range eh.msgBuffer[:eh.ptr] {
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

	_, err = eh.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("inbrief"),
		Key:    aws.String(fmt.Sprintf("%d.json", id)),
		Body:   bytes.NewReader(marshalled),
	})
	if err != nil {
		zap.L().Error("Failed to upload messages to S3", zap.Error(err))
		return err
	}

	zap.L().Info("Successfully flushed messages", zap.Int("count", eh.ptr), zap.String("id", fmt.Sprintf("%d", id)))

	eh.rdb.Publish(ctx, eh.pubChannel, id)
	eh.ptr = 0

	return nil
}

func (eh *EventHandler) Handle(ctx context.Context, listener *client.Listener, rdb *redis.Client) {
	for update := range listener.Updates {
		switch msg := update.(type) {
		case *client.UpdateNewMessage:
			eh.newMessageHandler(ctx, msg)
		}
	}
}

func (eh *EventHandler) FlushByPeriod(ctx context.Context, period time.Duration) {
	for {
		<-time.After(period)
		eh.flush(ctx)
	}
}

func processText(text *client.FormattedText) string {
	blacklist := []string{
		"textEntityTypeBotCommand",
		"textEntityTypeHashtag",
		"textEntityTypeMention",
		"textEntityTypeCashtag",
		"textEntityTypeMentionName",
		"textEntityTypeUrl",
		"textEntityTypeTextUrl",
	}

	sort.Slice(text.Entities, func(i, j int) bool {
		return text.Entities[i].Offset > text.Entities[j].Offset
	})

	u16Text := utf16.Encode([]rune(text.Text))

	for _, entity := range text.Entities {
		if slices.Contains(blacklist, entity.Type.TextEntityTypeType()) {
			if entity.Offset >= 0 &&
				entity.Offset < int32(len(u16Text)) &&
				entity.Offset+entity.Length <= int32(len(u16Text)) {

				u16Text = append(
					u16Text[:entity.Offset],
					u16Text[entity.Offset+entity.Length:]...)
			}
		}
	}

	return string(utf16.Decode(u16Text))
}

func (eh *EventHandler) newMessageHandler(ctx context.Context, msg *client.UpdateNewMessage) {
	zap.L().Debug(
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
		zap.L().Debug(
			"New message text",
			zap.String("text", content.Text.Text),
		)
		processedText := processText(content.Text)
		zap.L().Debug("Processed text", zap.String("text", processedText))
		if len([]rune(processedText)) > 50 {
			eh.msgBuffer[eh.ptr] = &pb.Message{
				Id:     msg.Message.Id,
				Text:   processedText,
				Ts:     timestamppb.New(time.Unix(int64(msg.Message.Date), 0)),
				ChatId: msg.Message.ChatId,
			}
			eh.ptr += 1

			if eh.ptr == len(eh.msgBuffer) {
				eh.flush(ctx)
			}
		}
	}
}
