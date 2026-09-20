package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	authjwt "kampiun/kernel/auth"
)

type ctxKey string

const principalKey ctxKey = "principal"

type Principal struct {
	UserID string
	Role   string
}

func MiddlewareAuth(tokens *authjwt.TokenMaker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			writeErr(w, http.StatusUnauthorized, "token hilang")
			return
		}
		claims, err := tokens.Parse(token)
		if err != nil {
			writeErr(w, http.StatusUnauthorized, "token tidak valid")
			return
		}
		ctx := context.WithValue(r.Context(), principalKey, Principal{UserID: claims.UserID, Role: claims.Role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PrincipalOK mengambil principal dari context.
func principal(r *http.Request) (Principal, bool) {
	p, ok := r.Context().Value(principalKey).(Principal)
	return p, ok
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
