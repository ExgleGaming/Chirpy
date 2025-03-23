package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	var chirpList []Chirp
	for _, chirp := range chirps {
		chirpList = append(chirpList, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirpList)
}
