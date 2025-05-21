package server

import (
	"context"
	"github.com/nrydanov/inbrief/config"
	"github.com/nrydanov/inbrief/gen/proto/fetcher"
	pc "github.com/nrydanov/inbrief/gen/proto/fetcher/fetcherconnect"
	"github.com/nrydanov/inbrief/internal"
	"github.com/nrydanov/inbrief/internal/tl"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/swaggest/swgui/v5emb"
	"github.com/zelenin/go-tdlib/client"
)

type server struct {
	state *internal.AppState
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

func StartServer(cfg *config.Config, state *internal.AppState) {
	path, handler := pc.NewFetcherServiceHandler(server{state: state})

	mux := http.NewServeMux()
	mux.Handle(path, handler)
	mux.HandleFunc("/api/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./gen/proto/fetcher/fetch.openapi.yaml")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.Handle("/api/docs/", v5emb.New(
		"Inbrief Scraper",
		"/api/swagger.yaml",
		"/api/docs/",
	))

	err := http.ListenAndServe(cfg.Server.GetAddr(), mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
