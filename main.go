package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nohlachilders/atlas/pkg/server"
)

func main() {
	ctx := context.Background()

	if err := server.Run(ctx, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "error in serving: %v\n", err)
		os.Exit(1)
	}
}
