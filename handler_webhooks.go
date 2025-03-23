package main

import (
	"context"
	"encoding/json"
	"github.com/exglegaming/Chirpy/internal/auth"
	"github.com/exglegaming/Chirpy/internal/database"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerUpdateUserChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}
	if key != cfg.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpdateUserChirpyRed(context.Background(), database.UpdateUserChirpyRedParams{
		ID:          params.Data.UserID,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User could not be found", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
