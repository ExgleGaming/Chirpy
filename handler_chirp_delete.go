package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/exglegaming/Chirpy/internal/auth"
	"github.com/exglegaming/Chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	chirpIDStr := parts[len(parts)-1]

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// First check if chirp exists
	chirp, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	// Then check if user is authorized
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps", nil)
		return
	}

	// Now perform the actual deletion
	_, err = cfg.db.DeleteChirp(context.Background(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userID,
	})
	if err != nil {
		// Try to determine what kind of error it is
		if strings.Contains(err.Error(), "not found") || errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		if strings.Contains(err.Error(), "not authorized") || strings.Contains(err.Error(), "permission") {
			respondWithError(w, http.StatusForbidden, "You can only delete your own chirps", err)
			return
		}
		// Generic error
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
