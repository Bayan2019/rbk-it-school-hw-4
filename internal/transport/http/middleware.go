package http

import (
	"net/http"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-4/internal/domain"
	// httptransport "github.com/Bayan2019/rbk-it-school-hw-4/internal/transport/http"
)

// 3. Middleware
func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Чтение Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				WriteError(w, http.StatusUnauthorized, "missing Authorization header", nil)
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				WriteError(w, http.StatusUnauthorized, "invalid Authorization header format", nil)
				return
			}

			// 2. Проверка токена
			tokenString := strings.TrimPrefix(authHeader, prefix)
			if tokenString == "" {
				WriteError(w, http.StatusUnauthorized, "empty token", nil)
				return
			}

			// 3. Парсинг claims
			claims, err := jwtManager.Validate(tokenString)
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "invalid token", err)
				return
			}

			// 4. Кладём user в context
			user := domain.UserContext{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			}

			ctx := withUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowedRoles ...domain.Roles) func(http.Handler) http.Handler {
	allowed := make(map[domain.Roles]struct{}, len(allowedRoles))
	// print()
	for _, role := range allowedRoles {
		// log.Println(role)
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := UserFromContext(r.Context())
			if err != nil {
				WriteError(w, http.StatusUnauthorized, "unauthorized", err)
				return
			}

			// log.Println(user)

			if _, ok := allowed[user.Role]; !ok {
				WriteError(w, http.StatusForbidden, "forbidden: insufficient role", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
