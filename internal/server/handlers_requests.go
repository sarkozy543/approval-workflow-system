package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sarkozy543/approval-workflow-system/internal/approval"
)


func (s *Server) handleListRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requests, err := s.approvalStore.GetAll(ctx)
	if err != nil {
		http.Error(w, "failed to fetch requests", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, requests)
}

func (s *Server) handleCreateRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var in approval.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	req, err := s.approvalStore.Create(ctx, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, req)
}
func (s *Server) handleGetRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// URL'den id parametresini al ( /requests/{id} )
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	req, err := s.approvalStore.GetByID(ctx, id)
	if err != nil {
		// Kayıt bulunamazsa 404
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}

		// Diğer hatalar 500
		http.Error(w, "failed to fetch request", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, req)
}

// approve/reject body için basit input
type actionInput struct {
	Note string `json:"note"`
}

func (s *Server) handleApproveRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Kullanıcıyı header'dan al (şimdilik basit)
	user := r.Header.Get("X-User")
	if user == "" {
		user = "anonymous"
	}

	var in actionInput
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&in) // note boş olabilir; hatayı önemsemiyoruz
	}

	req, err := s.approvalStore.Approve(ctx, id, user, in.Note)
	if err != nil {
		if errors.Is(err, approval.ErrInvalidStatus) {
        http.Error(w, "request is not in PENDING status", http.StatusBadRequest)
        return
    }
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to approve request", http.StatusInternalServerError)
		return
	}
	

	writeJSON(w, http.StatusOK, req)
}

func (s *Server) handleRejectRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	user := r.Header.Get("X-User")
	if user == "" {
		user = "anonymous"
	}

	var in actionInput
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&in)
	}

	req, err := s.approvalStore.Reject(ctx, id, user, in.Note)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "request not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to reject request", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, req)
}
func (s *Server) handleGetRequestLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	logs, err := s.approvalStore.GetLogsForRequest(ctx, id)
	if err != nil {
		http.Error(w, "failed to fetch logs", http.StatusInternalServerError)
		return
	}

	// Eğer hiç log yoksa, boş liste döneriz: []
	writeJSON(w, http.StatusOK, logs)
}

// Küçük yardımcı fonksiyon: JSON response yazmak için
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
