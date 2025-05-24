package server

import (
	"context"
	"fmt"

	"github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/nrydanov/inbrief/internal/tl"

	connect "connectrpc.com/connect"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
)

func (s server) Fetch(
	ctx context.Context,
	req *connect.Request[fetcher.FetchRequest],
) (*connect.Response[fetcher.FetchResponse], error) {
	state := s.state
	resp := &fetcher.FetchResponse{}
	info, err := state.TlClient.CheckChatFolderInviteLink(
		&client.CheckChatFolderInviteLinkRequest{
			InviteLink: req.Msg.ChatFolderLink,
		},
	)
	if err != nil {
		return nil, err
	}

	ids := tl.ExtractChatIds(info)

	zap.L().Debug("Scraping channels", zap.String("ids", fmt.Sprintf("%+v", ids)))

	for _, id := range ids {
		msgs, err := tl.FetchChannel(
			int64(id),
			req.Msg.LeftBound.AsTime(),
			req.Msg.RightBound.AsTime(),
			state,
		)
		// TODO(nrydanov): Handle error
		if err != nil {
			continue
		}

		resp.Messages = append(resp.Messages, msgs...)

	}

	go func() {
		for _, msg := range resp.Messages {
			s.msgCh <- msg
		}
		zap.L().Debug("All messages are sent to message channel")
	}()

	return connect.NewResponse(resp), nil
}

func (s server) SubscribeChat(
	ctx context.Context,
	req *connect.Request[fetcher.SubscribeChatFolderRequest],
) (*connect.Response[fetcher.Empty], error) {
	state := s.state

	info, err := state.TlClient.CheckChatFolderInviteLink(
		&client.CheckChatFolderInviteLinkRequest{
			InviteLink: req.Msg.ChatFolderLink,
		},
	)

	_ = tl.ExtractChatIds(info)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse[fetcher.Empty](nil), nil
}
