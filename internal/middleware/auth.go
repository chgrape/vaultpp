package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/chgrape/vaultpp/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jwtKey := os.Getenv("JWT_SECRET")

		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Empty Authorization header", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "Invalid authorization header scheme", http.StatusBadRequest)
			return
		}

		splitHeader := strings.Split(header, " ")

		token, err := jwt.ParseWithClaims(splitHeader[1], &service.Claims{}, func(t *jwt.Token) (any, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
