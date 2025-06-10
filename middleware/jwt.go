package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"ph.certs.com/clm_main/auth"
)

// verifyToken validates a JWT token string and returns an error if it's invalid or cannot be parsed.
func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return auth.SecretKey, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

// JWTMiddleware should always check for the presence of the jwt token in the Authorization header
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
