package security

import (
	"encoding/json"
	"maps"
	"net/http"
	"strings"
	"time"
)

type tencentCloudTokenExchangeRequest struct {
	SecretID         string   `json:"secret_id"`
	SecretKey        string   `json:"secret_key"`
	SessionToken     string   `json:"session_token"`
	DisplayName      string   `json:"display_name"`
	AllowedRegions   []string `json:"allowed_regions"`
	Scopes           []string `json:"scopes"`
	ExpiresInSeconds int64    `json:"expires_in_seconds"`
	Description      string   `json:"description"`
}

type TokenExchangeBootstrapConfig struct {
	MCPServerName string
	ServerURL     string
}

func RegisterTokenExchangeRoutes(mux *http.ServeMux, service *TencentCloudTokenExchangeService, bootstrapCfg TokenExchangeBootstrapConfig) {
	if mux == nil || service == nil {
		return
	}
	mux.Handle("/auth/token-exchange/tencentcloud", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleTencentCloudTokenExchange(w, r, service, bootstrapCfg, false)
	}))
	mux.Handle("/auth/bootstrap/tencentcloud", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleTencentCloudTokenExchange(w, r, service, bootstrapCfg, true)
	}))
}

func handleTencentCloudTokenExchange(w http.ResponseWriter, r *http.Request, service *TencentCloudTokenExchangeService, bootstrapCfg TokenExchangeBootstrapConfig, includeBootstrap bool) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	var req tencentCloudTokenExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
		return
	}
	result, err := service.Exchange(r.Context(), TencentCloudTokenExchangeInput{
		SecretID:        req.SecretID,
		SecretKey:       req.SecretKey,
		SessionToken:    req.SessionToken,
		DisplayName:     req.DisplayName,
		AllowedRegions:  req.AllowedRegions,
		RequestedScopes: req.Scopes,
		ExpiresIn:       time.Duration(req.ExpiresInSeconds) * time.Second,
		Description:     req.Description,
	})
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "failed") || strings.Contains(strings.ToLower(err.Error()), "empty") {
			status = http.StatusUnauthorized
		}
		writeJSON(w, status, map[string]any{"error": err.Error()})
		return
	}
	response := tokenExchangeResponse(result)
	if includeBootstrap {
		mergeResponse(response, buildBootstrapResponse(bootstrapCfg, result.Issued.Token))
	}
	writeJSON(w, http.StatusCreated, response)
}

func tokenExchangeResponse(result *TencentCloudTokenExchangeResult) map[string]any {
	response := map[string]any{
		"id":              result.Issued.Record.ID,
		"token":           result.Issued.Token,
		"token_prefix":    result.Issued.Record.TokenPrefix,
		"subject_type":    result.Issued.Record.SubjectType,
		"subject_id":      result.Issued.Record.SubjectID,
		"tenant_id":       result.Issued.Record.TenantID,
		"display_name":    result.Issued.Record.DisplayName,
		"scopes":          result.Issued.Record.Scopes,
		"allowed_regions": result.Issued.Record.AllowedRegions,
		"expires_at":      result.Issued.Record.ExpiresAt.Format(time.RFC3339),
		"credential_kind": result.CredentialKind,
		"identity": map[string]any{
			"type":         result.Identity.Type,
			"account_id":   result.Identity.AccountID,
			"user_id":      result.Identity.UserID,
			"principal_id": result.Identity.PrincipalID,
			"arn":          result.Identity.Arn,
			"request_id":   result.Identity.RequestID,
		},
	}
	if result.CredentialExpiresAt != nil {
		response["credential_expires_at"] = result.CredentialExpiresAt.Format(time.RFC3339)
	}
	return response
}

func buildBootstrapResponse(cfg TokenExchangeBootstrapConfig, token string) map[string]any {
	serverName := strings.TrimSpace(cfg.MCPServerName)
	if serverName == "" {
		serverName = "mcp-server-postgres"
	}
	serverURL := strings.TrimSpace(cfg.ServerURL)
	mcpServerConfig := map[string]any{
		"type": "sse",
		"url":  serverURL,
	}
	if token != "" {
		mcpServerConfig["headers"] = map[string]string{
			"Authorization": "Bearer " + token,
		}
	}
	mcpConfig := map[string]any{
		"mcpServers": map[string]any{
			serverName: mcpServerConfig,
		},
	}
	prettyJSON, _ := json.MarshalIndent(mcpConfig, "", " ")
	response := map[string]any{
		"mcp_server_name": serverName,
		"mcp_server_url":  serverURL,
		"mcp_config":      mcpConfig,
		"mcp_config_json": string(prettyJSON),
		"client_env": map[string]string{
			"MCP_ACCESS_TOKEN": token,
			"MCP_SERVER_URL":   serverURL,
		},
		"next_steps": []string{
			"Copy mcp_config_json into your MCP client configuration.",
			"Or export MCP_ACCESS_TOKEN and point your client to mcp_server_url.",
		},
	}
	return response
}

func mergeResponse(dst map[string]any, extra map[string]any) {
	maps.Copy(dst, extra)
}
