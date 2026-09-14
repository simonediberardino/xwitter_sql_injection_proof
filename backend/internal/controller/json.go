package controller

import (
	"encoding/json"
	"net/http"

	"xwitter/internal/dto"
)

func decodeJSON(r *http.Request, dest any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 16*1024)
	return json.NewDecoder(r.Body).Decode(dest)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, dto.ErrorResponse{Error: message})
}
