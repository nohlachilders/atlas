package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
)

func Run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	mux := http.ServeMux{}
	server := http.Server{
		Addr:    ":8080",
		Handler: &mux,
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	return server.ListenAndServe()
}

type Config struct {
	Port string
}
