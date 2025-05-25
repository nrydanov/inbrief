package tl

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"time"
	"unicode/utf16"

	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/redis/go-redis/v9"
	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventHandler struct {
	listener *client.Listener
	outputCh chan<- *pb.Message
	client   *client.Client
}

func NewEventHandler(
	outputCh chan<- *pb.Message,
	client *client.Client,
	bufferSize int,
) *EventHandler {
	return &EventHandler{
		client:   client,
		listener: client.GetListener(),
		outputCh: outputCh,
	}
}

func (eh *EventHandler) Handle(
	ctx context.Context,
	listener *client.Listener,
	rdb *redis.Client,
) {
	for {
		select {
		case update := <-listener.Updates:
			switch msg := update.(type) {
			case *client.UpdateNewMessage:
				err := eh.newMessageHandler(msg)
				if err != nil {
					zap.L().Error("Unable to handle new message", zap.Error(err))
				}
				continue
			}
		case <-ctx.Done():
			return
		}
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

func (eh *EventHandler) newMessageHandler(msg *client.UpdateNewMessage) error {
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

		chat, err := eh.client.GetChat(&client.GetChatRequest{
			ChatId: msg.Message.ChatId,
		})
		if err != nil {
			zap.L().Error("Unable to get chat")
			return err
		}

		username, err := ExtractUsername(eh.client, chat)
		if err != nil {
			zap.L().Error("Unable to extract username", zap.Error(err))
			return err
		}

		if len([]rune(processedText)) > 50 {
			eh.outputCh <- &pb.Message{
				Text: processedText,
				Ts:   timestamppb.New(time.Unix(int64(msg.Message.Date), 0)),
				Link: fmt.Sprintf("https://t.me/%s/%d", username, msg.Message.Id),
			}
			zap.L().Debug("Processed text is sent to output channel")
		}
	}

	return nil
}
