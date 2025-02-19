package server

import (
	"context"
	"fmt"
	"net/http"
	//"testing"
	"time"
)

func awaitServerStartup(ctx context.Context, endpoint string, timeout time.Duration) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err.Error())
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error making request: %s", err.Error())
		}

		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			fmt.Println("server is ready")
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) > timeout {
				return fmt.Errorf("timeout reached")
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
}
