package main

import (
	"net/http"

	"github.com/MarunDArbaumont/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error while retreiving refresh token", err)
		return
	}

	err = cfg.database.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "refresh token not found or expired", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}