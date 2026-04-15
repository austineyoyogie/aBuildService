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

// get authenticated user in angular 19 example
// Get the current login user 17,18
// https://www.google.com/search?q=get+authenticated+user+in+angular+20&sca_esv=d0849dc9e7ce7768&ei=GkJwacOIH-Wm0PEPyIyliAg&ved=0ahUKEwjD_aPl0puSAxVlEzQIHUhGCYEQ4dUDCBE&uact=5&oq=get+authenticated+user+in+angular+20&gs_lp=Egxnd3Mtd2l6LXNlcnAiJGdldCBhdXRoZW50aWNhdGVkIHVzZXIgaW4gYW5ndWxhciAyMDIFECEYoAEyBRAhGKABMgUQIRigATIFECEYoAEyBRAhGKABSP4qUPwLWNkgcAF4AJABAJgBYaAB_AOqAQE3uAEDyAEA-AEBmAIHoAL5A8ICCBAAGLADGO8FwgILEAAYsAMYogQYiQXCAgUQIRirAsICBRAhGJ8FmAMAiAYBkAYFkgcDNi4xoAfKJrIHAzUuMbgH8gPCBwUwLjMuNMgHF4AIAA&sclient=gws-wiz-serp
// https://www.google.com/search?q=get+authenticated+user+in+angular+20&sca_esv=d0849dc9e7ce7768&ei=GkJwacOIH-Wm0PEPyIyliAg&ved=0ahUKEwjD_aPl0puSAxVlEzQIHUhGCYEQ4dUDCBE&uact=5&oq=get+authenticated+user+in+angular+20&gs_lp=Egxnd3Mtd2l6LXNlcnAiJGdldCBhdXRoZW50aWNhdGVkIHVzZXIgaW4gYW5ndWxhciAyMDIFECEYoAEyBRAhGKABMgUQIRigATIFECEYoAEyBRAhGKABSP4qUPwLWNkgcAF4AJABAJgBYaAB_AOqAQE3uAEDyAEA-AEBmAIHoAL5A8ICCBAAGLADGO8FwgILEAAYsAMYogQYiQXCAgUQIRirAsICBRAhGJ8FmAMAiAYBkAYFkgcDNi4xoAfKJrIHAzUuMbgH8gPCBwUwLjMuNMgHF4AIAA&sclient=gws-wiz-serp
// https://www.google.com/search?q=get+authenticated+user+in+angular+18&sca_esv=d0849dc9e7ce7768&ei=jj5wacWDPZWd0PEPzIvwuAM&ved=0ahUKEwjFgsa0z5uSAxWVDjQIHcwFHDcQ4dUDCBE&uact=5&oq=get+authenticated+user+in+angular+18&gs_lp=Egxnd3Mtd2l6LXNlcnAiJGdldCBhdXRoZW50aWNhdGVkIHVzZXIgaW4gYW5ndWxhciAxODIFECEYoAEyBRAhGKABMgUQIRigATIFECEYoAEyBRAhGKABSJ0WUMYGWL8LcAF4AZABAJgBYaABtgGqAQEyuAEDyAEA-AEBmAIDoALJAcICChAAGLADGNYEGEfCAgUQIRirApgDAIgGAZAGCJIHAzIuMaAH9AqyBwMxLjG4B8EBwgcDMi0zyAcOgAgA&sclient=gws-wiz-serp
// https://www.google.com/search?q=get+authenticated+user+in+angular+17&sca_esv=d0849dc9e7ce7768&source=hp&ei=Vz5waeD9LKeOur8Pt9Df4Qc&iflsig=AFdpzrgAAAAAaXBMZ38SZRlwZtBuHLgUAI3pYQR3sIFT&oq=get+authenticated+user+in+angu&gs_lp=Egdnd3Mtd2l6Ih5nZXQgYXV0aGVudGljYXRlZCB1c2VyIGluIGFuZ3UqAggDMgUQIRigATIFECEYoAEyBRAhGKABMgUQIRigATIFECEYoAEyBRAhGJ8FMgUQIRifBUiCpgNQ9A9YzvgCcBB4AJABAJgBW6AB5hWqAQI0N7gBA8gBAPgBAZgCP6AC6heoAgrCAh0QABiABBi0AhjUAxjlAhjnBhi3AxiKBRjqAhiKA8ICCxAAGIAEGJECGIoFwgILEAAYgAQYsQMYgwHCAhEQLhiABBixAxjRAxiDARjHAcICEBAuGIAEGNEDGEMYxwEYigXCAgUQLhiABMICBRAAGIAEwgIOEAAYgAQYsQMYgwEYigXCAggQLhiABBixA8ICCBAAGIAEGLEDwgIOEC4YgAQYsQMYxwEYrwHCAgsQLhiABBjHARivAcICCxAuGIAEGNEDGMcBwgILEC4YgAQYsQMYigXCAgsQLhiABBixAxiDAcICDBAAGIAEGLEDGAoYC8ICCxAAGIAEGIYDGIoFwgIFEAAY7wXCAgcQABiABBgNwgIGEAAYDRgewgIKEAAYBRgKGA0YHsICChAAGAgYChgNGB7CAggQABiiBBiJBcICCBAAGIAEGKIEwgIGEAAYFhgewgIIEAAYFhgKGB7CAgcQIRigARgKwgIFECEYqwKYAwbxBc2FJlOd4IdikgcCNjOgB6nWArIHAjQ3uAerF8IHBzAuMTMuNTDIB7oBgAgA&sclient=gws-wiz
// https://www.google.com/search?q=How+to+get+authenticated+user+in+golang+api+web+page&sca_esv=3e001843e4d244d0&biw=1792&bih=983&aic=0&ei=4QNuafa_AcTE0PEPrYmF4QY&ved=0ahUKEwi2tMGVr5eSAxVEIjQIHa1EIWwQ4dUDCBM&uact=5&oq=How+to+get+authenticated+user+in+golang+api+web+page&gs_lp=Egxnd3Mtd2l6LXNlcnAiNEhvdyB0byBnZXQgYXV0aGVudGljYXRlZCB1c2VyIGluIGdvbGFuZyBhcGkgd2ViIHBhZ2UyBRAhGKsCSP5RUIALWNFLcAF4AJABAJgBdaABngaqAQQxMC4xuAEDyAEA-AEBmAILoALTBsICBRAhGKABwgIFEAAY7wXCAggQABiABBiiBJgDAIgGAZIHBDEwLjGgB4YqsgcEMTAuMbgH0wbCBwUwLjUuNsgHIIAIAQ&sclient=gws-wiz-serp
// https://www.google.com/search?sca_esv=3e001843e4d244d0&q=How+to+get+authenticated+user+in+golang+api+web+page+example+qui&sa=X&ved=2ahUKEwii8c7yuZeSAxWxJEQIHb01CXgQ1QJ6BAg6EAE&biw=1781&bih=930&dpr=2&aic=0

/*
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
*/
