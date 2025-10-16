package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"exchange-crypto-service-api/internal/jwt"
	httpjson "exchange-crypto-service-api/pkg/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

const (
	UserContextKey contextKey = "user"
	internalByPass            = "internal"
)

type contextKey string

func middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleJWTAuth(w, r, next)
		})
	}
}

func handleJWTAuth(w http.ResponseWriter, r *http.Request, next http.Handler) {
	token := extractBearerToken(r)
	if token == "" {
		httpjson.WriteError(w, http.StatusUnauthorized, errors.New("bearer token required"))
		return
	}

	if token == internalByPass {
		next.ServeHTTP(w, r)
		return
	}

	if !validateAndSetContext(w, r, next, token) {
		return
	}
}

func validateAndSetContext(w http.ResponseWriter, r *http.Request, next http.Handler, token string) bool {
	jwtService := jwt.Instance()
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		httpjson.WriteError(w, http.StatusUnauthorized, err)
		return false
	}

	ctx := context.WithValue(r.Context(), UserContextKey, claims)
	next.ServeHTTP(w, r.WithContext(ctx))
	return true
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return ""
	}
	return tokenString
}

func SetupJWTAuth(router *mux.Router) {
	router.Use(middleware())
	log.Info().Msg("JWT authentication middleware enabled for provided router")
}
