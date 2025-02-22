package server

import (
	"context"
	"testing"
)

func TestStartupWait(t *testing.T) {
	ctx := context.Background()
	testingStartDefault(t, ctx)
}
