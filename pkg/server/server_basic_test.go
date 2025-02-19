package server

import (
	"context"
	"testing"
	"time"
)

func TestStartupWait(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	baseURL := "http://localhost" + defaultTestingGetenv("PORT")

	go Run(ctx, defaultTestingGetenv)

	if err := awaitServerStartup(ctx, baseURL+"/healthz", 5*time.Second); err != nil {
		t.Errorf("error in awaiting server startup: %v", err)
	}
}

func defaultTestingGetenv(s string) string {
	envs := map[string]string{
		"PORT": ":8080",
	}
	if value, ok := envs[s]; ok {
		return value
	}
	return ""
}
