package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neuroshepherd/learn-http-servers/internal/auth"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	expiresIn := time.Hour

	token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	validatedUserID, err := auth.ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("Expected userID %v, got %v", userID, validatedUserID)
	}
}

// test an already expired token
func TestJWTExpired(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	expiresIn := -time.Hour

	token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = auth.ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Expected error validating expired JWT, got nil")
	}
}

// create a JWT with one secret, then try to validate it with a different secret and expect an error
func TestJWTInvalidSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	invalidSecret := "my-invalid-secret-key"
	expiresIn := time.Hour

	token, err := auth.MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = auth.ValidateJWT(token, invalidSecret)
	if err == nil {
		t.Fatal("Expected error validating JWT with invalid secret, got nil")
	}
}
