package server

import (
	"dsc/inbrief/scraper/config"
	"dsc/inbrief/scraper/internal"
	"log"
	"net/http"

	pb "dsc/inbrief/scraper/pkg/proto"
	"github.com/swaggest/swgui/v5emb"
)

type server struct {
	state *internal.AppState
}

func StartServer(cfg *config.Config, state *internal.AppState) {
	twirpHandler := pb.NewFetcherServer(&server{state: state})

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	mux.HandleFunc("/api/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/openapi.json")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.Handle("/api/docs/", v5emb.New(
		"Inbrief Scraper",
		"/api/swagger.json",
		"/api/docs/",
	))

	err := http.ListenAndServe(cfg.Server.GetAddr(), mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
