package main

import (
	"net/http"
	"time"

	"github.com/MarunDArbaumont/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error while retreiving refresh token", err)
		return
	}

	refreshTokenDB, err := cfg.database.GetRefreshTokeByToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token not found or expired", err)
		return
	}


	token, err := auth.MakeJWT(refreshTokenDB.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while creating JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: token,
	})
}