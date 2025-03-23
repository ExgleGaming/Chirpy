package main

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	if author == "" {
		if sort == "desc" {
			chirps, err := cfg.db.GetChirpsDesc(context.Background())
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
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
			return
		}
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
		return
	} else {
		user, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't parse author_id", err)
			return
		}
		if sort == "desc" {
			chirps, err := cfg.db.GetChirpsByUserID(context.Background(), user)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't find chirps by user", err)
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
			return
		} else {
			chirps, err := cfg.db.GetChirpsByUserID(context.Background(), user)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't find chirps by user", err)
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
			return
		}
	}
}
