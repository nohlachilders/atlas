package server

import (
	"context"
	"testing"
)

func TestUserCreation(t *testing.T) {
	ctx := context.Background()
	baseURL, ctx := testingStartDefault(t, ctx)

	testResetDatabase(t, ctx, baseURL)

	status, body, err := testingSendRequestWithJSON(ctx, baseURL+"/users", "POST", map[string]string{"email": "test@email.com", "password": "test"})
	if status != "201 Created" {
		t.Errorf("something went wrong with the user creation endpoint: wrong status %s", status)
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	t.Logf("body: %s", string(body))
}

func testingUserCreation(t *testing.T, ctx context.Context, baseURL string, data any) (string, string) {
	status, body, err := testingSendRequestWithJSON(ctx, baseURL+"/users", "POST", data)
	t.Logf("%v,%v,%v", status, string(body), err)
	if err != nil {
		t.Errorf("error in creation: %s", err.Error())
	}
	if status != "201 Created" {
		t.Errorf("error in creation: %s", body)
	}
	return status, string(body)
}
