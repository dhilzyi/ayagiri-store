package middleware

import (
	"log"
	"net/http"
	"time"

	"restaurant/internal/database"

	"github.com/goccy/go-json"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Middleware struct {
	db *database.Queries
}

func NewMiddleware(db *database.Queries) *Middleware {
	return &Middleware{
		db: db,
	}
}

func Chain(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func Logging(next http.Handler) http.Handler {
	counter := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		log.Printf("%s %s counter: %d", r.Method, r.URL.Path, counter)

		next.ServeHTTP(w, r)
	})
}

func Test(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}

		session, err := m.db.GetSession(r.Context(), cookie.Value)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		if session.ExpiresAt.Time.Before(time.Now()) {
			if err := m.db.DeleteSession(r.Context(), cookie.Value); err != nil {
			}
			respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
