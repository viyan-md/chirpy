package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

const testSecret = "super-secret-dev-key"

func TestJWT_HappyPath(t *testing.T) {
	uid := uuid.New()

	tok, err := MakeJWT(uid, testSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned err: %v", err)
	}

	got, err := ValidateJWT(tok, testSecret)
	if err != nil {
		t.Fatalf("ValidateJWT returned err: %v", err)
	}
	if got != uid {
		t.Fatalf("want %v, got %v", uid, got)
	}
}

func TestJWT_Expired(t *testing.T) {
	uid := uuid.New()

	tok, _ := MakeJWT(uid, testSecret, -time.Minute)

	if _, err := ValidateJWT(tok, testSecret); err == nil {
		t.Fatalf("expected error for expired token, got nil")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	uid := uuid.New()

	tok, _ := MakeJWT(uid, testSecret, time.Hour)

	if _, err := ValidateJWT(tok, "wrong-secret"); err == nil {
		t.Fatalf("expected error for wrong secret, got nil")
	}
}
