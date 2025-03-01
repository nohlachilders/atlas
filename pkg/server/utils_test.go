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

// testingSendRequestWithJSON sends a request to the given endpoint with the given data,
// and returns the status, the raw data of the response body, and any errors.
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

// testingSendAuthedRequestWithJSON sends a request to the given endpoint with the given data and auth token,
// and returns the status, the raw data of the response body, and any errors.
func testingSendAuthedRequestWithJSON(ctx context.Context, endpoint string, requestType string, data any, token string) (status string, response []byte, err error) {
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
	req.Header.Add("Authorization", "Bearer "+token)

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
	return status, response, nil
}

// testingLogin is a helper function to test the login endpoint. returns the status,
// and the body of the response.
func testingRefresh(t *testing.T, ctx context.Context, baseURL string, token string) (string, string) {
	buffer := bytes.NewReader([]byte{})
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/refresh", buffer)
	if err != nil {
		t.Errorf("error in creation: %s", err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("error in creation: %s", err.Error())
	}

	status := resp.Status
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("error in creation: %s", err.Error())
	}

	resp.Body.Close()
	data := map[string]any{}
	json.Unmarshal([]byte(response), &data)
	token, ok := data["token"].(string)
	if !ok {
		t.Errorf("malformed response: %v", data)
	}

	return status, token
}

// testingResetDatabase is a helper function to reset the database.
// this endpoint only works if the server is configured to the "dev" platform.
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
		"ATLAS_SECRET":   "testsecret",
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

// testingAwaitServerStartup waits for the server to be ready by making a request to the given endpoint.
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
