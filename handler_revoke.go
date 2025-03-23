package main

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Getting token if it exists in the header
	authHeader := r.Header.Get("Authorization")
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		http.Error(w, "malformed authorization header", http.StatusUnauthorized)
		return
	}

	// Getting the refresh token
	refreshToken, err := cfg.db.GetRefreshTokenByToken(context.Background(), splitAuth[1])
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token is not valid", err)
		return
	}

	// Check if token is expired or revoked
	if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token is expired or revoked", nil)
		return
	}

	err = cfg.db.UpdateRefreshToken(context.Background(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
