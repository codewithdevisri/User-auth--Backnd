package authorization

import (
	"context"
	"github.com/Ayan25844/netflix/token"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextKeyCompanyId contextKey = "CompanyId"
	ContextKeyUsername  contextKey = "Username"
	ContextKeyRoles     contextKey = "Roles"
)

// ValidateToken middleware - compatible with net/http
func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the "Authorization" header in Bearer format
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized - invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		referer := r.Header.Get("Referer")
		valid, claims := token.VerifyToken(tokenString, referer)
		if !valid {
			http.Error(w, "Unauthorized - invalid or missing token", http.StatusUnauthorized)
			return
		}

		// Add claims to context and proceed to next handler
		ctx := context.WithValue(r.Context(), ContextKeyCompanyId, claims.ID)
		ctx = context.WithValue(ctx, ContextKeyUsername, claims.Name)
		ctx = context.WithValue(ctx, ContextKeyRoles, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Authorization(validRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rolesVal := r.Context().Value(ContextKeyRoles)
			if rolesVal == nil {
				http.Error(w, "Unauthorized - roles not found", http.StatusUnauthorized)
				return
			}

			roles, ok := rolesVal.([]string)
			if !ok {
				http.Error(w, "Unauthorized - invalid roles type", http.StatusUnauthorized)
				return
			}

			// Convert slice to map for O(1) lookups
			roleMap := make(map[string]struct{}, len(roles))
			for _, role := range roles {
				roleMap[strings.ToUpper(role)] = struct{}{}
			}

			// Check if user has at least one of the valid roles
			authorized := false
			for _, validRole := range validRoles {
				if _, found := roleMap[strings.ToUpper(validRole)]; found {
					authorized = true
					break
				}
			}

			if !authorized {
				http.Error(w, "Unauthorized - insufficient role", http.StatusUnauthorized)
				return
			}

			// Passed authorization, proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}
