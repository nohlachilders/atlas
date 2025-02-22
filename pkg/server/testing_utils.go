package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	//"testing"
	"time"
)

func testingStartDefault(t *testing.T, ctx context.Context) {
	env := testingEnv{
		"ATLAS_PORT":   ":8080",
		"ATLAS_DB_URL": "postgresql://localhost:5432/atlas?sslmode=disable",
	}

	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	baseURL := "http://localhost" + env.getEnv("ATLAS_PORT")

	go Run(ctx, env.getEnv)
	time.Sleep(10 * time.Millisecond)

	if err := testingAwaitServerStartup(ctx, baseURL+"/healthz", 5*time.Second); err != nil {
		t.Errorf("error in awaiting server startup: %v", err)
	}
}

func testingAwaitServerStartup(ctx context.Context, endpoint string, timeout time.Duration) error {
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

type testingEnv map[string]string

func (env testingEnv) getEnv(s string) string {
	if value, ok := env[s]; ok {
		return value
	}
	return ""
}
