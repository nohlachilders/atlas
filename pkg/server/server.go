package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/nohlachilders/atlas/internal/database"
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
	cfg.Context = &ctx

	cfg.Port = getenv("ATLAS_PORT")
	if cfg.Port == "" {
		fmt.Printf("PORT is empty. Defaulting to :8080\n")
		cfg.Port = ":8080"
	}

	dbURL := getenv("ATLAS_DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil || db.Ping() != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	cfg.Database = database.New(db)

	cfg.Platform = getenv("ATLAS_PLATFORM")

	mux := http.ServeMux{}
	cfg.makeRoutes(&mux)
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
	Port     string
	Platform string
	Database *database.Queries
	Context  *context.Context
}
