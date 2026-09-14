package controller

import (
	"net/http"
	"strings"

	"xwitter/internal/auth"
	"xwitter/internal/dto"
)

func (c *Controller) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.posts.List(r.Context(), dto.PostFilter{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"posts": posts})
}

func (c *Controller) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "Post cannot be empty")
		return
	}
	if len(req.Content) > 280 {
		writeError(w, http.StatusBadRequest, "Post must be 280 characters or less")
		return
	}

	claims := auth.UserFromRequest(r)
	created, err := c.posts.Create(r.Context(), dto.Post{
		UserID:  claims.ID,
		Content: req.Content,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"post": created})
}
