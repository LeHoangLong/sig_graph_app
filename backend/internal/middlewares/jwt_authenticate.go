package middlewares

import (
	"backend/internal/controllers"
	"backend/internal/services"
	"net/http"
)

type JwtAuthenticateMiddleware struct {
	userController *controllers.UserController
}

func (m *JwtAuthenticateMiddleware) authenticate(iCookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				for _, val := range r.Cookies() {
					if val.Name == iCookieName {
						token = val.Value
					}
				}

			}

			// Allow unauthenticated users in
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			username, err := m.userController.LogInWithToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// put it in context
			ctx := services.PutUsernameInContex(r.Context(), username)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
