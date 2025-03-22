package auth

import (
	"net/http"
	"strings"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		expectErr bool
	}{
		{
			name:      "Valid Bearer Token",
			headers:   http.Header{"Authorization": []string{"Bearer my-secret-token"}},
			wantToken: "my-secret-token",
			expectErr: false,
		},
		{
			name:      "Missing Authorization Header",
			headers:   http.Header{},
			wantToken: "",
			expectErr: true,
		},
		{
			name:      "Invalid Authorization Format",
			headers:   http.Header{"Authorization": []string{"MyToken secret-token"}},
			wantToken: "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got error: %v", tt.expectErr, err)
			}
			if strings.TrimSpace(gotToken) != tt.wantToken {
				t.Errorf("expected token: %q, got token: %q", tt.wantToken, gotToken)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
