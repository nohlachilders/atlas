package server

import (
	"context"
	"testing"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	baseURL, ctx := testingStartDefault(t, ctx)

	testResetDatabase(t, ctx, baseURL)

	cases := []map[string]any{
		{
			"email":    "john@beatles.com",
			"password": "john1",
		},
		{
			"email":    "paul@beatles.com",
			"password": "paul2",
		},
		{
			"email":    "george@beatles.com",
			"password": "george3",
		},
		{
			"email":    "ringo@beatles.com",
			"password": "ringo4",
		},
	}

	for _, c := range cases {
		_, _ = testingUserCreation(t, ctx, baseURL, c)

		_, _ = testingLogin(t, ctx, baseURL, c)
	}
}

func testingLogin(t *testing.T, ctx context.Context, baseURL string, data any) (string, string) {
	status, body, err := testingSendRequestWithJSON(ctx, baseURL+"/login", "POST", data)
	t.Logf("%v,%v,%v", status, string(body), err)
	if err != nil {
		t.Errorf("error in creation: %s", err.Error())
	}
	if status != "200 OK" {
		t.Errorf("error in creation: %s", body)
	}
	return status, string(body)
}
