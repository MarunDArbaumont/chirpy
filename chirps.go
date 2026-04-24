package main

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/MarunDArbaumont/chirpy/internal/database"
	"github.com/MarunDArbaumont/chirpy/internal/auth"
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
	}
	type returnVals struct {
		Chirp
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while retreiving token", err)
		return
	}

	currentUserID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "log in before posting chirp", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "chirp is too long", nil)
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
    	UserID: currentUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong when creating chirp", err)
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
		respondWithError(w, http.StatusInternalServerError, "something went wrong when retrieving chirps", err)
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
		respondWithError(w, http.StatusBadRequest, "the id given is not an id", err)
		return
	}
	chirp, err := cfg.database.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found or doesn't exist anymore", err)
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