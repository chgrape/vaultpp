package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chgrape/vaultpp/internal/service"
	"github.com/chgrape/vaultpp/internal/vault"
	"github.com/golang-jwt/jwt/v5"
)

type VaultProvider struct {
	JwtKey string
}

func NewVaultProvider(s *vault.VaultService) (*VaultProvider, error) {
	jwtKey, err := s.FetchSecret("auth", "jwt_key")
	if err != nil {
		return nil, err
	}
	return &VaultProvider{JwtKey: jwtKey}, nil
}

type key string

const ClaimsCtxKey key = "claims"

func (p *VaultProvider) AuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			return []byte(p.JwtKey), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*service.Claims)

		ctx := context.WithValue(context.Background(), ClaimsCtxKey, *claims)

		handler(w, r.WithContext(ctx))
	}
}
