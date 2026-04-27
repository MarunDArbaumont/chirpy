package main

import (
	"net/http"
	"encoding/json"
	"time"
	"database/sql"

	"github.com/google/uuid"
	"github.com/MarunDArbaumont/chirpy/internal/database"
	"github.com/MarunDArbaumont/chirpy/internal/auth"
)

type User struct {
	ID        		uuid.UUID `json:"id"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedAt 		time.Time `json:"updated_at"`
	Email     		string    `json:"email"`
	Password  		string	`json:"hashed_password"`
	Token 	  		string 	`json:"token"`
	RefreshToken 	string	`json:"refresh_token"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	 string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
	}

	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		Email:   params.Email,
    	HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "email already exist or not valid", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		User: User {
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
}


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	 	string 	`json:"email"`
		Password 	string 	`json:"password"`
	}
	type returnVals struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	user, err := cfg.database.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "this email is not valid", err)
		return
	}

	isSame, err := auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't compare passwords", err)
		return
	}

	if !isSame {
		respondWithError(w, http.StatusUnauthorized, "the password is incorrect", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while creating JWT", err)
		return
	}

	refreshToken, err := cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: auth.MakeRefreshToken(),
		UserID: user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while creating refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: User {
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: token,
			RefreshToken:  refreshToken.Token,
		},
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	 	string 	`json:"email"`
		Password 	string 	`json:"password"`
	}

	type returnVals struct {
		User
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error while retreiving token", err)
		return
	}

	currentUserID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "log in before posting chirp", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	updateUser, err := cfg.database.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
		ID: currentUserID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "user not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: User {
			ID: updateUser.ID,
			CreatedAt: updateUser.CreatedAt,
			UpdatedAt: updateUser.UpdatedAt,
			Email: updateUser.Email,
		},
	})
}
