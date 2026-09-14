package controller

import "net/http"

func (c *Controller) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/auth/register", c.Register)
	mux.HandleFunc("POST /api/auth/login", c.Login)
	mux.HandleFunc("GET /api/auth/me", c.requireAuth(c.Me))
	mux.HandleFunc("GET /api/posts", c.ListPosts)
	mux.HandleFunc("POST /api/posts", c.requireAuth(c.CreatePost))
	mux.HandleFunc("POST /api/users/search", c.SearchUsers)
	mux.HandleFunc("GET /api/users/{username}", c.GetUser)
	mux.HandleFunc("GET /api/users/{username}/posts", c.GetUserPosts)
	return mux
}
