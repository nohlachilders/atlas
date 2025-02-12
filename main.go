package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nohlachilders/atlas/internal/server"
)

func main() {
	ctx := context.Background()

	fmt.Println("serving...")
	if err := server.Run(ctx, []string{}, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "error in serving: %v\n", err)
		os.Exit(1)
	}
}
