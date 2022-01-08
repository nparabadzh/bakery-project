package rest

import (
	"bakery-project/entities"
	"context"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const BEARER_SCHEMA = "Bearer"
		//var header = r.Header.Get("x-access-token") //Grab the token from the header
		var header = r.Header.Get("Authorization") //Grab the token from the header

		var token string
		if header == "" || len(header) <= len(BEARER_SCHEMA) {
			respondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}
		token = header[len(BEARER_SCHEMA):]
		token = strings.TrimSpace(token)

		if token == "" {
			//Token is missing, returns with error code 401 Unauthorized
			respondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}
		claims := &entities.UserToken{}

		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
