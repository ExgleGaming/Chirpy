package main

import (
	"context"
	"github.com/exglegaming/Chirpy/internal/auth"
	"net/http"
	"strings"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

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
	
	// Make JWT for user
	token, err := auth.MakeJWT(refreshToken.UserID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
