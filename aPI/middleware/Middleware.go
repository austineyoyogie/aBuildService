package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type User struct {
	ID      string
	Email   string
	IsAdmin bool
}

type ContextKey string

const (
	ClaimsContextKey ContextKey = "userID"
)

type AuthKey struct{}

func AuthMiddlewareFunc(tokenMaker *JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error verifying token: %v", err), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(ClaimsContextKey).(int)
	if !ok {
		return -1
	}
	return userID
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value(ClaimsContextKey).(uuid.UUID)
	return userID, ok
}

func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *JWTMaker) (*UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}
	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}
	token := fields[1]
	claims, err := tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return claims, nil
}
func Logger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: Method %s, Path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
func DetectBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "firefox") {
		return "Mozilla Firefox"
	} else if strings.Contains(ua, "chrome") {
		return "Google Chrome"
	} else if strings.Contains(ua, "safari") {
		return "Apple Safari"
	} else if strings.Contains(ua, "edge") {
		return "Microsoft Edge"
	}
	return "Unknown Browser"
}

func GetAuthMiddlewareFunc(tokenMaker *JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token.......: %v", err), http.StatusUnauthorized)
				return
			}
			//if !claims.IsAdmin {
			//	http.Error(w, "user is not an admin", http.StatusForbidden)
			//	return
			//}
			ctx := context.WithValue(r.Context(), AuthKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func RequestAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

