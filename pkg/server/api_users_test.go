package server

import (
	"context"
	"encoding/json"
	"strings"
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
	if !strings.Contains(string(body), "test") {
		t.Errorf("unexpected missing text: %s", string(body))
	}
	//t.Logf("body: %s", string(body))
}

func TestUserInfo(t *testing.T) {
	ctx := context.Background()
	baseURL, ctx := testingStartDefault(t, ctx)

	testResetDatabase(t, ctx, baseURL)

	cases := map[string]map[string]any{
		"john": {
			"email":    "john@beatles.com",
			"password": "john1",
			"token":    "",
		},
		"paul": {
			"email":    "paul@beatles.com",
			"password": "paul2",
			"token":    "",
		},
		"george": {
			"email":    "george@beatles.com",
			"password": "george3",
			"token":    "",
		},
		"ringo": {
			"email":    "ringo@beatles.com",
			"password": "ringo4",
			"token":    "",
		},
	}

	for _, c := range cases {
		_, _ = testingUserCreation(t, ctx, baseURL, c)

		_, resp := testingLogin(t, ctx, baseURL, c)
		data := map[string]any{}
		json.Unmarshal([]byte(resp), &data)

		refresh := data["refresh_token"].(string)
		_, token := testingRefresh(t, ctx, baseURL, refresh)
		c["token"] = token

		status, body := testingFetchUserInfo(t, ctx, baseURL, token)
		if status != "200 OK" {
			t.Errorf("unexpected status: %v\n%v", status, body)
		}
	}

	//t.Errorf("%v", cases["john"]["token"])
	token := cases["john"]["token"].(string)
	status, body := testingUpdateUserInfo(t, ctx, baseURL, token, map[string]any{"email": "lennon@beatles.com", "password": "lennon1"})
	if status != "200 OK" {
		t.Errorf("unexpected status: %v\n%v", status, body)
	}
	if !strings.Contains(body, "lennon") {
		t.Errorf("user not updated properly: %v\n%v", status, body)
	}
	status, body = testingDeleteUser(t, ctx, baseURL, token)
	if status != "200 OK" {
		t.Errorf("unexpected status: %v\n%v", status, body)
	}

	status, body = testingFetchUserInfo(t, ctx, baseURL, token)
	if status != "500 Internal Server Error" {
		t.Errorf("unexpected status: %v\n%v", status, body)
	}

}

// testingFetchUserInfo is a helper function to test the user info endpoint.
// returns the status, and the body of the response.
func testingFetchUserInfo(t *testing.T, ctx context.Context, baseURL string, token string) (string, string) {
	status, body, err := testingSendAuthedRequestWithJSON(ctx, baseURL+"/users", "GET", nil, token)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	return status, string(body)
}

// testingUpdateUserInfo is a helper function to test the user update endpoint.
// returns the status, and the body of the response.
func testingUpdateUserInfo(t *testing.T, ctx context.Context, baseURL string, token string, data any) (string, string) {
	status, body, err := testingSendAuthedRequestWithJSON(ctx, baseURL+"/users", "PUT", data, token)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	if status != "200 OK" {
		t.Errorf("error: %s", err.Error())
	}
	return status, string(body)
}

// testingDeleteUser is a helper function to test the user deletion endpoint.
// returns the status, and the body of the response.
func testingDeleteUser(t *testing.T, ctx context.Context, baseURL string, token string) (string, string) {
	status, body, err := testingSendAuthedRequestWithJSON(ctx, baseURL+"/users", "DELETE", nil, token)
	if err != nil {
		t.Errorf("error: %s", err.Error())
	}
	if status != "200 OK" {
		t.Errorf("error: %s", err.Error())
	}
	return status, string(body)
}

// testingUserCreation is a helper function to test the user creation endpoint.
// returns the status, and the body of the response.
func testingUserCreation(t *testing.T, ctx context.Context, baseURL string, data any) (string, string) {
	status, body, err := testingSendRequestWithJSON(ctx, baseURL+"/users", "POST", data)
	//t.Logf("%v,%v,%v", status, string(body), err)
	if err != nil {
		t.Errorf("error in creation: %s", err.Error())
	}
	if status != "201 Created" {
		t.Errorf("error in creation: %s", body)
	}
	return status, string(body)
}
