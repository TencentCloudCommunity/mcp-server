package security

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type createTokenRequest struct {
	SubjectType      string   `json:"subject_type"`
	SubjectID        string   `json:"subject_id"`
	TenantID         string   `json:"tenant_id"`
	DisplayName      string   `json:"display_name"`
	Scopes           []string `json:"scopes"`
	AllowedRegions   []string `json:"allowed_regions"`
	ExpiresInSeconds int64    `json:"expires_in_seconds"`
	Description      string   `json:"description"`
}

type revokeTokenRequest struct {
	Reason string `json:"reason"`
}

func RegisterAdminRoutes(mux *http.ServeMux, issuer *TokenIssuer, store TokenStore) {
	if mux == nil || issuer == nil || store == nil {
		return
	}
	mux.Handle("/admin/tokens", WrapAdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateToken(w, r, issuer)
		case http.MethodGet:
			handleListTokens(w, r, store)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		}
	})))
	mux.Handle("/admin/tokens/", WrapAdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remainder := strings.TrimPrefix(r.URL.Path, "/admin/tokens/")
		if remainder == "" {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		if id, ok := strings.CutSuffix(remainder, "/revoke"); ok {
			id = strings.TrimSuffix(id, "/")
			if r.Method != http.MethodPost || strings.Contains(id, "/") || id == "" {
				writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
				return
			}
			handleRevokeToken(w, r, store, id)
			return
		}
		if r.Method != http.MethodGet || strings.Contains(remainder, "/") {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		handleGetToken(w, r, store, remainder)
	})))
}

func handleCreateToken(w http.ResponseWriter, r *http.Request, issuer *TokenIssuer) {
	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
		return
	}
	issued, err := issuer.Issue(r.Context(), CreateTokenInput{
		SubjectType:    req.SubjectType,
		SubjectID:      req.SubjectID,
		TenantID:       req.TenantID,
		DisplayName:    req.DisplayName,
		Scopes:         req.Scopes,
		AllowedRegions: req.AllowedRegions,
		ExpiresIn:      time.Duration(req.ExpiresInSeconds) * time.Second,
		Description:    req.Description,
		IssuedBy:       adminActorFromRequest(r),
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"id":              issued.Record.ID,
		"token":           issued.Token,
		"token_prefix":    issued.Record.TokenPrefix,
		"subject_type":    issued.Record.SubjectType,
		"subject_id":      issued.Record.SubjectID,
		"tenant_id":       issued.Record.TenantID,
		"display_name":    issued.Record.DisplayName,
		"scopes":          issued.Record.Scopes,
		"allowed_regions": issued.Record.AllowedRegions,
		"status":          issued.Record.Status,
		"expires_at":      issued.Record.ExpiresAt.Format(time.RFC3339),
		"created_at":      issued.Record.CreatedAt.Format(time.RFC3339),
	})
}

func handleListTokens(w http.ResponseWriter, r *http.Request, store TokenStore) {
	records, err := store.ListTokens(r.Context(), TokenFilter{
		TenantID:  strings.TrimSpace(r.URL.Query().Get("tenant_id")),
		SubjectID: strings.TrimSpace(r.URL.Query().Get("subject_id")),
		Status:    TokenStatus(strings.TrimSpace(r.URL.Query().Get("status"))),
		Scope:     strings.TrimSpace(r.URL.Query().Get("scope")),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "list tokens failed"})
		return
	}
	items := make([]map[string]any, 0, len(records))
	for _, record := range records {
		items = append(items, tokenRecordResponse(record))
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items, "total": len(items)})
}

func handleGetToken(w http.ResponseWriter, r *http.Request, store TokenStore, id string) {
	record, err := store.GetTokenByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "token not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "get token failed"})
		return
	}
	writeJSON(w, http.StatusOK, tokenRecordResponse(*record))
}

func handleRevokeToken(w http.ResponseWriter, r *http.Request, store TokenStore, id string) {
	var req revokeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
		return
	}
	if err := store.RevokeToken(r.Context(), id, RevokeTokenInput{RevokedBy: adminActorFromRequest(r), Reason: req.Reason}); err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "token not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "revoke token failed"})
		return
	}
	if credentialStore, ok := store.(CredentialStore); ok {
		if err := credentialStore.DeleteCredentialBinding(r.Context(), id); err != nil && !errors.Is(err, ErrCredentialBindingNotFound) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "revoke token succeeded but credential cleanup failed"})
			return
		}
	}
	record, err := store.GetTokenByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "status": TokenStatusRevoked})
		return
	}
	writeJSON(w, http.StatusOK, tokenRecordResponse(*record))
}

func adminActorFromRequest(r *http.Request) string {
	if r == nil {
		return "admin"
	}
	if actor := strings.TrimSpace(r.Header.Get(AdminActorHeader)); actor != "" {
		return actor
	}
	return "admin"
}

func tokenRecordResponse(record TokenRecord) map[string]any {
	response := map[string]any{
		"id":              record.ID,
		"token_prefix":    record.TokenPrefix,
		"subject_type":    record.SubjectType,
		"subject_id":      record.SubjectID,
		"tenant_id":       record.TenantID,
		"display_name":    record.DisplayName,
		"scopes":          record.Scopes,
		"allowed_regions": record.AllowedRegions,
		"status":          record.Status,
		"expires_at":      record.ExpiresAt.Format(time.RFC3339),
		"created_at":      record.CreatedAt.Format(time.RFC3339),
		"updated_at":      record.UpdatedAt.Format(time.RFC3339),
		"issued_by":       record.IssuedBy,
		"description":     record.Description,
	}
	if record.LastUsedAt != nil {
		response["last_used_at"] = record.LastUsedAt.Format(time.RFC3339)
	}
	if record.RevokedBy != "" {
		response["revoked_by"] = record.RevokedBy
	}
	if record.RevokedAt != nil {
		response["revoked_at"] = record.RevokedAt.Format(time.RFC3339)
	}
	if record.RevokeReason != "" {
		response["revoke_reason"] = record.RevokeReason
	}
	return response
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
