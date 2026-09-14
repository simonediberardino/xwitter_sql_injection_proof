package controller

import (
	"net/http"
	"strings"

	"xwitter/internal/auth"
	"xwitter/internal/repository"
)

type Controller struct {
	users *repository.UserRepository
	posts *repository.PostRepository
	auth  *auth.Service
}

func New(users *repository.UserRepository, posts *repository.PostRepository, authService *auth.Service) *Controller {
	return &Controller{users: users, posts: posts, auth: authService}
}

func (c *Controller) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		claims, err := c.auth.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		next(w, r.WithContext(auth.WithUser(r.Context(), claims)))
	}
}

func CORS(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
