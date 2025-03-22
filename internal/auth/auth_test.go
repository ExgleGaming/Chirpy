package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"

	// Test 1: Create a token and validate it successfully
	t.Run("Valid JWT", func(t *testing.T) {
		// Create a token that expires after an hour
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("Error creating JWT: %v", err)
		}

		// Validate token
		extractedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("Error validating JWT: %v", err)
		}

		// Check if extractes ID matches the original
		if extractedID != userID {
			t.Errorf("Expected user ID %v, go %v", userID, extractedID)
		}
	})

	// Test 2: Check that expired tokens are rejected
	t.Run("Expired JWT", func(t *testing.T) {
		// Create a token that expired 1 sec ago
		token, err := MakeJWT(userID, secret, -1*time.Second)
		if err != nil {
			t.Fatalf("Error creating JWT: %v", err)
		}

		// Try to validate the expired token
		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("Expected error for expired token, got nil")
		}
	})

	// Test 3: Check that tokens with wrong secret are rejected
	t.Run("Wrong Secret", func(t *testing.T) {
		// Create a token
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("Error creating JWT: %v", err)
		}

		// Validating that it will catch a wrong key being passed through
		_, err = ValidateJWT(token, "wrong-secret-key")
		if err == nil {
			t.Fatalf("Expected error for token with wrong secret, got nil")
		}
	})

	// Test 4: Invalid token format
	t.Run("Invalid Token Format", func(t *testing.T) {
		// Try to validate a malformed token
		invalidToken := "not.a.valid.jwt.token"
		_, err := ValidateJWT(invalidToken, secret)
		if err == nil {
			t.Error("Expected error for invalid token format, got nil")
		}
	})
}
