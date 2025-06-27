package middleware

import (
	"context"
	"log"
	"net/http"
	"shop-backend/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const (
	UserContextKey  key = "userEmail"
	AdminContextKey key = "adminEmail"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Attach user email to context
		email := claims["email"].(string)
		ctx := context.WithValue(r.Context(), UserContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// tokenStr := extractToken(r)
		// if tokenStr == "" {
		// http.Error(w, "Missing token", http.StatusUnauthorized)
		// return
		// }

		tokenStr := ""

		if cookie, err := r.Cookie("token"); err == nil {
			tokenStr = cookie.Value
		} else {
			tokenStr = extractToken(r) // fallback to header
		}
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(tokenStr)
		log.Println("This is the admin middleware")
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		log.Println("This the admin middleware at end")

		cfg := config.LoadConfig()
		email := claims["email"].(string)
		if email != cfg.AdminEmail {
			http.Error(w, "Unauthorized admin access", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), AdminContextKey, email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Extract token from Authorization header: Bearer <token>
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

// Validate JWT token and return claims
func validateToken(tokenStr string) (jwt.MapClaims, error) {

	cfg := config.LoadConfig()
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	log.Print("Tracking......................,", err)
	log.Println(token.Valid)
	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
