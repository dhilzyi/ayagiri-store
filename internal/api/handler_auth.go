package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"restaurant/internal/auth"
	"restaurant/internal/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const (
	EXPIRE_TIME = time.Hour * 24
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var payloadLogin LoginRequest
	if err := decodeJson(r.Body, &payloadLogin); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	user, err := h.db.GetUserByUsername(ctx, payloadLogin.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "invalid credentials", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	isValid, err := auth.VerifyHashPassword(payloadLogin.Password, user.PasswordHash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}
	if !isValid {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", fmt.Errorf("invalid password"))
		return
	}

	token, err := auth.GenerateSessionToken()
	if err != nil {
		respondWithError(w, http.StatusFailedDependency, "", err)
		return
	}

	session, err := h.db.CreateSession(ctx, database.CreateSessionParams{
		UserID:    pgtype.Int4{Valid: true, Int32: user.ID},
		ExpiresAt: pgtype.Timestamp{Valid: true, Time: time.Now().Add(EXPIRE_TIME)},
		Token:     token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  time.Now().Add(EXPIRE_TIME),
		HttpOnly: true,
		Path:     "/",
	})

	respondWithJSON(w, http.StatusCreated, nil)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if err := h.db.DeleteSession(r.Context(), cookie.Value); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete session id", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	session, err := auth.ValidateSession(r, h.db)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if err := h.db.UpdateExpireSession(r.Context(), database.UpdateExpireSessionParams{
		ExpiresAt: pgtype.Timestamp{Valid: true, Time: session.ExpiresAt.Time.Add(EXPIRE_TIME)},
		Token:     session.Token,
	}); err != nil {
		log.Printf("failed to update session expired time: %s\n", session.Token)
	}

	respondWithJSON(w, http.StatusOK, session.Token)
}
