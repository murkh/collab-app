package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/murkh/collab-app/collab-app/internal/store"
)

type Signer interface {
	SignToken(claims jwt.MapClaims) (string, error)
}

type Handler struct {
	store  *store.Store
	signer Signer
}

func NewHandler(s *store.Store, signer Signer) *Handler {
	return &Handler{store: s, signer: signer}
}

type tokenReq struct {
	DocID string `json:"docId"`
}

func (h *Handler) IssueCollabToken(w http.ResponseWriter, r *http.Request) {
	uid := r.Header.Get("X-User-ID")
	if uid == "" {
		http.Error(w, "missing user header", http.StatusUnauthorized)
	}

	var req tokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	docUUID, err := uuid.Parse(req.DocID)
	if err != nil {
		http.Error(w, "invalid doc id", http.StatusBadRequest)
		return
	}

	ok, role, err := h.store.CanUserAccess(r.Context(), docUUID, uid)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	exp := time.Now().Add(15 * time.Minute).Unix()
	claims := jwt.MapClaims{
		"sub":  uid,
		"doc":  req.DocID,
		"role": role,
		"exp":  exp,
	}
	signed, err := h.signer.SignToken(claims)
	if err != nil {
		http.Error(w, "server error signing token", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{"token": signed})

}
