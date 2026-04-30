package main

import (
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/MarunDArbaumont/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event 	 string `json:"event"`
		Data struct{
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	APIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "you are not allowed to do this action, no API key given", err)
		return
	}

	if cfg.polkaKey != APIKey {
		respondWithError(w, http.StatusUnauthorized, "your API key doesn't match", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.database.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}