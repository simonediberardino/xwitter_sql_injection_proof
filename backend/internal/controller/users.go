package controller

import (
	"errors"
	"net/http"
	"strings"

	"xwitter/internal/db"
	"xwitter/internal/dto"
)

func (c *Controller) SearchUsers(w http.ResponseWriter, r *http.Request) {
	var query dto.SearchQuery
	if err := decodeJSON(r, &query); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	query.Query = strings.TrimSpace(query.Query)
	if query.Query == "" {
		writeJSON(w, http.StatusOK, []dto.User{})
		return
	}

	rows, err := c.users.Search(r.Context(), query)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, rows)
}

func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := c.users.FindByUsername(r.Context(), dto.User{Username: r.PathValue("username")})
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user, err = c.users.CountPosts(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (c *Controller) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	user, err := c.users.FindByUsername(r.Context(), dto.User{Username: r.PathValue("username")})
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	posts, err := c.posts.List(r.Context(), dto.PostFilter{UserID: user.ID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"posts": posts})
}
