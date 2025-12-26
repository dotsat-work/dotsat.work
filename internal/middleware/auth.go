package middleware

import (
	"net/http"

	"dotsat.work/internal/ctxkeys"
	"dotsat.work/internal/service"
	"github.com/google/uuid"
)

// AuthMiddleware checks for JWT token and adds user + tenant to context if valid
func AuthMiddleware(authService *service.AuthService, userService *service.UserService, tenantService *service.TenantService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get JWT from a cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				// No cookie, continue without auth
				next.ServeHTTP(w, r)
				return
			}

			// Verify token
			claims, err := authService.VerifyJWT(cookie.Value)
			if err != nil {
				// Invalid token, clear the cookie and continue
				authService.ClearJWTCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			// Get user ID from claims
			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				authService.ClearJWTCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				authService.ClearJWTCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			// Fetch user from the database
			user, err := userService.ByID(userID)
			if err != nil {
				authService.ClearJWTCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			// Security: Remove password hash from context
			user.PasswordHash = nil

			// Fetch tenant
			tenant, err := tenantService.ByID(user.TenantID)
			if err != nil {
				authService.ClearJWTCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			// Add user + tenant to context
			ctx := ctxkeys.WithUser(r.Context(), user)
			ctx = ctxkeys.WithTenant(ctx, tenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth ensures the user is authenticated
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := ctxkeys.User(r.Context())
		if user == nil {
			// For HTMX requests, use HX-Redirect header to force full page redirect
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/auth")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			// For regular requests, use standard redirect
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		// TODO: Add onboarding check when ProfileService is implemented
		// profile := ctxkeys.Profile(r.Context())
		// if profile != nil && profile.Name == "" && r.URL.Path != "/auth/onboarding" {
		//     redirect to /auth/onboarding
		// }

		next.ServeHTTP(w, r)
	}
}

// RequireGuest ensures the user is not authenticated
func RequireGuest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := ctxkeys.User(r.Context())
		if user != nil {
			// For HTMX requests, use HX-Redirect header to force full page redirect
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/app/dashboard")
				w.WriteHeader(http.StatusSeeOther)
				return
			}
			// For regular requests, use standard redirect
			http.Redirect(w, r, "/app/dashboard", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
