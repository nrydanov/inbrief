package server

import (
	"context"
	pb "github.com/nrydanov/inbrief/gen/proto/fetcher"
	"github.com/nrydanov/inbrief/internal/tl"
	"fmt"

	connect "connectrpc.com/connect"

	"github.com/zelenin/go-tdlib/client"
	"go.uber.org/zap"
)

func (s server) Fetch(
	ctx context.Context,
	req *connect.Request[pb.FetchRequest],
) (*connect.Response[pb.FetchResponse], error) {
	state := s.state
	resp := &pb.FetchResponse{}
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

	return connect.NewResponse(resp), nil
}
