package server

import (
	"context"
	"testing"
)

func TestUserCreation(t *testing.T) {
	ctx := context.Background()
	testingStartDefault(t, ctx)
}
