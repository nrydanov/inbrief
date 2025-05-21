package server

import (
	"dsc/inbrief/scraper/config"
	pc "dsc/inbrief/scraper/gen/proto/fetcher/protoconnect"
	"dsc/inbrief/scraper/internal"
	"github.com/swaggest/swgui/v5emb"
	"log"
	"net/http"
)

type server struct {
	state *internal.AppState
}

func StartServer(cfg *config.Config, state *internal.AppState) {
	path, handler := pc.NewFetcherServiceHandler(server{state: state})

	mux := http.NewServeMux()
	mux.Handle(path, handler)
	mux.HandleFunc("/api/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/proto/fetch.openapi.yaml")
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
