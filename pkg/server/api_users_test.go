package server

import (
	"context"
	"testing"
)

func TestUserCreation(t *testing.T) {
	ctx := context.Background()
	baseURL, ctx := testingStartDefault(t, ctx)
	status, _, err := testingSendRequestWithJSON(ctx, baseURL+"/reset", "POST", "")
	if status != "204 No Content" {
		t.Errorf("something went wrong with the reset endpoint: wrong status %s", status)
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	status, body, err := testingSendRequestWithJSON(ctx, baseURL+"/users", "POST", map[string]string{"email": "test@email.com", "password": "test"})
	if status != "201 Created" {
		t.Errorf("something went wrong with the user creation endpoint: wrong status %s", status)
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	t.Logf("body: %s", string(body))

	testingSendRequestWithJSON(ctx, baseURL+"/reset", "POST", "")
}
