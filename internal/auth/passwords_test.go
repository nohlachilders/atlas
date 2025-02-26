package auth

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	pass := "test"
	_, err := HashPassword(pass)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestCheckPassword(t *testing.T) {
	pass := "test"
	hashed := "$2a$10$5CFy00EAVTrbdRNbLHR0z.dnTbGp1X9E.5STOeKn95h/kcjjAFlcC"
	err := CheckPasswordHash(pass, hashed)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
