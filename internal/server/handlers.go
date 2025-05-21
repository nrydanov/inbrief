package server

import (
	connect "connectrpc.com/connect"
	"context"
	pb "dsc/inbrief/scraper/gen/proto"
	"dsc/inbrief/scraper/pkg/tl"
	"fmt"
	"time"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s server) Fetch(
	ctx context.Context,
	req *connect.Request[pb.FetchRequest],
) (resp *connect.Response[pb.FetchResponse], err error) {
	state := s.state
	fetchResponse := &pb.FetchResponse{}

	info, err := state.TlClient.CheckChatFolderInviteLink(
		&client.CheckChatFolderInviteLinkRequest{
			InviteLink: req.Msg.ChatFolderLink,
		},
	)
	if err != nil {
		return nil, err
	}

	ids := make([]tl.ChatId, len(info.AddedChatIds))

	for i, id := range info.AddedChatIds {
		ids[i] = tl.ChatId(id)
	}

	zap.L().Debug("Scraping channels", zap.String("ids", fmt.Sprintf("%+v", ids)))

	for _, id := range ids {
		fromMessageId := int64(0)
		for {
			history, err := state.TlClient.GetChatHistory(
				&client.GetChatHistoryRequest{
					ChatId:        int64(id),
					FromMessageId: fromMessageId,
					Limit:         100,
				},
			)
			if err != nil {
				return nil, err
			}

			reachedEnd := false

			for _, message := range history.Messages {
				if int64(message.Date) < req.Msg.LeftBound.Seconds {
					zap.L().Debug("Reached left bound")
					reachedEnd = true
					break
				}
				zap.L().Debug("Scraped message", zap.String("time", fmt.Sprintf("%+v", message.Date)))

				switch message.Content.(type) {
				case *client.MessageText:
					fetchResponse.Messages = append(fetchResponse.Messages, &pb.Message{
						Id:     message.Id,
						Text:   message.Content.(*client.MessageText).Text.Text,
						Ts:     timestamppb.New(time.Unix(int64(message.Date), 0)),
						ChatId: message.ChatId,
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
	}

	resp = connect.NewResponse(fetchResponse)

	return resp, err
}
