package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/MarunDArbaumont/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID 	  uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type returnVals struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	params.Body = replaceBadWord(params.Body, badWords)
	
	chirp, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
    	UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong when creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		Chirp: Chirp {
			ID: 		chirp.ID,
			CreatedAt: 	chirp.CreatedAt,
			UpdatedAt: 	chirp.UpdatedAt,
			Body: 		chirp.Body,
			UserID: 	chirp.UserID,
		},
	})
}

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.database.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong when retrieving chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: 		dbChirp.ID,
			CreatedAt: 	dbChirp.CreatedAt,
			UpdatedAt: 	dbChirp.UpdatedAt,
			Body: 		dbChirp.Body,
			UserID: 	dbChirp.UserID,
		})
	}


	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerSingleChirp(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Chirp
	}

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "The id given is not an id", err)
	}
	chirp, err := cfg.database.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found or doesn't exist anymore", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Chirp: Chirp {
			ID: 		chirp.ID,
			CreatedAt: 	chirp.CreatedAt,
			UpdatedAt: 	chirp.UpdatedAt,
			Body: 		chirp.Body,
			UserID: 	chirp.UserID,
		},
	})
}