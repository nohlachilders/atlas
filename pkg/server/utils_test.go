package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	//"testing"
	"time"
)

func testingSendRequestWithJSON(ctx context.Context, endpoint string, requestType string, data any) (status string, response []byte, err error) {
	client := http.Client{}
	body, err := json.Marshal(data)
	if err != nil {
		return "", nil, err
	}
	buffer := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, requestType, endpoint, buffer)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}

	status = resp.Status
	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	resp.Body.Close()
	return status, response, err
}

func testResetDatabase(t *testing.T, ctx context.Context, baseURL string) {
	status, body, err := testingSendRequestWithJSON(ctx, baseURL+"/reset", "POST", "")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if status != "204 No Content" {
		t.Errorf("something went wrong with the reset endpoint: wrong status %s, %s", status, body)
	}
}

// Start a server instance using the current test case and a given context.
// waits for the healthz enpoint to be ready before returning.
// returns the base url used by the server and the context generated for it
func testingStartDefault(t *testing.T, ctx context.Context) (string, context.Context) {
	env := Env{
		"ATLAS_PORT":     ":8080",
		"ATLAS_DB_URL":   "postgresql://localhost:5432/atlas?sslmode=disable",
		"ATLAS_PLATFORM": "dev",
	}
	getEnv := makeGetEnv(env)

	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)

	go Run(ctx, getEnv)

	time.Sleep(10 * time.Millisecond)

	baseURL := "http://localhost" + env["ATLAS_PORT"]
	if err := testingAwaitServerStartup(ctx, baseURL+"/healthz", 5*time.Second); err != nil {
		t.Errorf("error in awaiting server startup: %v", err)
	}
	return baseURL, ctx
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

type Env map[string]string

func makeGetEnv(env Env) func(string) string {
	return func(s string) string {
		if value, ok := env[s]; ok {
			return value
		}
		return ""
	}
}
