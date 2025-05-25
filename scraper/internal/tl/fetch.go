package tl

import (
	"fmt"
	"time"

	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/nrydanov/inbrief/internal"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FetchChannel(
	chId int64,
	leftBound time.Time,
	rightBound time.Time,
	state *internal.AppState,
) ([]*pb.Message, error) {
	messages := make([]*pb.Message, 0)

	// TODO(nrydanov): Add right bound support
	fromMessageId := int64(0)
	for {
		history, err := state.TlClient.GetChatHistory(
			&client.GetChatHistoryRequest{
				ChatId:        int64(chId),
				FromMessageId: fromMessageId,
				Limit:         100,
			},
		)
		if err != nil {
			zap.L().Debug("Unable to get chat history")
			return nil, err
		}

		chat, err := state.TlClient.GetChat(&client.GetChatRequest{
			ChatId: chId,
		})

		if err != nil {
			zap.L().Debug("Unable to get chat")
			return nil, err
		}

		username, err := ExtractUsername(state.TlClient, chat)
		if err != nil {
			zap.L().Debug("Unable to extract username", zap.Error(err))
			continue
		}

		zap.L().Debug("Chat info",
			zap.Any("chat", chat),
			zap.String("type", chat.Type.ChatTypeType()),
		)

		reachedEnd := false

		for _, message := range history.Messages {
			if int64(message.Date) < leftBound.Unix() {
				zap.L().Debug("Reached left bound")
				reachedEnd = true
				break
			}
			zap.L().Debug(
				"Scraped message",
				zap.String("time", fmt.Sprintf("%+v", message.Date)),
			)

			switch message.Content.(type) {
			case *client.MessageText:
				messages = append(messages, &pb.Message{
					Text: processText(message.Content.(*client.MessageText).Text),
					Ts:   timestamppb.New(time.Unix(int64(message.Date), 0)),
					Link: fmt.Sprintf("https://t.me/%s/%d", username, message.Id),
				})
			default:
				continue
			}
		}
		if reachedEnd {
			break
		}

		fromMessageId = history.Messages[len(history.Messages)-1].Id

	}

	return messages, nil
}
