package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
)

func Run(
	ctx context.Context,
	// we abstract the way we source environment variables to enable
	// in-code testing of the whole server process
	getenv func(string) string,
) error {
	cfg := Config{}
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg.Port = getenv("PORT")
	if cfg.Port == "" {
		fmt.Printf("PORT is empty. Defaulting to :8080\n")
		cfg.Port = ":8080"
	}

	mux := http.ServeMux{}
	makeRoutes(&mux)
	server := http.Server{
		Addr:    cfg.Port,
		Handler: &mux,
	}

	fmt.Println("Now serving...")
	serverError := make(chan (error))
	go func() { serverError <- server.ListenAndServe() }()
	for {
		select {
		case err := <-serverError:
			return err
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return nil
		}
	}
}

type Config struct {
	Port string
}
