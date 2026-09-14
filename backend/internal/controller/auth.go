package controller

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"xwitter/internal/auth"
	"xwitter/internal/db"
	"xwitter/internal/dto"
)

var (
	usernameRE = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	emailRE    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.Bio = strings.TrimSpace(req.Bio)

	switch {
	case !usernameRE.MatchString(req.Username):
		writeError(w, http.StatusBadRequest, "Username must be 3-20 characters and use only letters, numbers, or underscores")
		return
	case !emailRE.MatchString(req.Email):
		writeError(w, http.StatusBadRequest, "Enter a valid email address")
		return
	case len(req.DisplayName) < 2 || len(req.DisplayName) > 40:
		writeError(w, http.StatusBadRequest, "Display name must be 2-40 characters")
		return
	case len(req.Password) < 8:
		writeError(w, http.StatusBadRequest, "Password must be at least 8 characters")
		return
	case len(req.Bio) > 160:
		writeError(w, http.StatusBadRequest, "Bio must be 160 characters or less")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	created, err := c.users.Create(r.Context(), dto.User{
		Username:     req.Username,
		Email:        req.Email,
		DisplayName:  req.DisplayName,
		Bio:          req.Bio,
		PasswordHash: string(hash),
	})
	if err != nil {
		if errors.Is(err, db.ErrDuplicate) {
			writeError(w, http.StatusConflict, "Username or email is already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	token, err := c.auth.Sign(created.ID, created.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, dto.AuthResponse{Token: token, User: created})
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Identifier == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Username/email and password are required")
		return
	}

	user, err := c.users.FindByIdentifier(r.Context(), req)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := c.auth.Sign(user.ID, user.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.AuthResponse{Token: token, User: user})
}

func (c *Controller) Me(w http.ResponseWriter, r *http.Request) {
	claims := auth.UserFromRequest(r)
	user, err := c.users.FindByID(r.Context(), dto.User{ID: claims.ID})
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "Account no longer exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}
