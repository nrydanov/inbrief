package server

import (
	"context"
	"net/http"

	"github.com/nrydanov/inbrief/config"
	"github.com/nrydanov/inbrief/gen/proto/fetcher"
	pc "github.com/nrydanov/inbrief/gen/proto/fetcher/fetcherconnect"
	"github.com/nrydanov/inbrief/internal"
	"go.uber.org/zap"

	"github.com/swaggest/swgui/v5emb"
)

type server struct {
	state *internal.AppState
	msgCh chan *fetcher.Message
}

func StartServer(
	ctx context.Context,
	cfg *config.Config,
	state *internal.AppState,
	msgCh chan *fetcher.Message,
) {
	path, handler := pc.NewFetcherServiceHandler(server{
		state: state,
		msgCh: msgCh,
	})

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

	server := &http.Server{
		Addr:    cfg.Server.GetAddr(),
		Handler: mux,
	}

	go func() {
		if err := http.ListenAndServe(cfg.Server.GetAddr(), mux); err != http.ErrServerClosed {
			zap.L().Fatal("failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		zap.L().Fatal("failed to shutdown server", zap.Error(err))
	}

}
