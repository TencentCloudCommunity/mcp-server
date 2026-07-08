package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"postgres_server/security"
)

type transportMode string

const (
	transportModeStreamableHTTP transportMode = "streamable-http"
	transportModeSSE            transportMode = "sse"
	transportModeStdio          transportMode = "stdio"
	transportEnv                              = "MCP_TRANSPORT"
	defaultSTSRegion                          = "ap-guangzhou"
)

type transportConfig struct {
	Mode            transportMode
	BindHost        string
	Port            string
	HTTPEndpoint    string
	SSEEndpoint     string
	MessageEndpoint string
	StatelessHTTP   bool
	PublicURL       string
}

func transportConfigFromEnv() transportConfig {
	return transportConfig{
		Mode:            normalizeTransport(os.Getenv(transportEnv)),
		BindHost:        getEnv("MCP_SERVER_BIND_HOST", "127.0.0.1"),
		Port:            getEnv("MCP_SERVER_PORT", getEnv("MCP_SERVER_SSE_PORT", "9000")),
		HTTPEndpoint:    normalizeEndpointPath(getEnv("MCP_SERVER_HTTP_ENDPOINT", "/mcp")),
		SSEEndpoint:     normalizeEndpointPath(getEnv("MCP_SERVER_SSE_ENDPOINT", "/sse")),
		MessageEndpoint: normalizeEndpointPath(getEnv("MCP_SERVER_MESSAGE_ENDPOINT", "/message")),
		StatelessHTTP:   getEnvBool("MCP_STREAMABLE_HTTP_STATELESS", true),
		PublicURL:       strings.TrimSpace(os.Getenv("MCP_SERVER_PUBLIC_URL")),
	}
}

func (c transportConfig) UsesHTTP() bool {
	return c.Mode != transportModeStdio
}

func (c transportConfig) Validate(authMode security.AuthMode) error {
	if c.Mode == transportModeStdio && authMode == security.AuthModeIssuedToken {
		return fmt.Errorf("MCP_TRANSPORT=stdio does not support MCP_AUTH_MODE=issued-token on the current main branch; please use request-credential/shared-token/none, or switch transport to streamable-http/sse")
	}
	return nil
}

func (c transportConfig) ListenAddr() string {
	return fmt.Sprintf("%s:%s", c.BindHost, c.Port)
}

func (c transportConfig) ClientType() string {
	if c.Mode == transportModeSSE {
		return "sse"
	}
	return "streamable-http"
}

func (c transportConfig) ActiveEndpoint() string {
	switch c.Mode {
	case transportModeSSE:
		return c.SSEEndpoint
	case transportModeStdio:
		return ""
	default:
		return c.HTTPEndpoint
	}
}

func (c transportConfig) ServerURL() string {
	if c.Mode == transportModeStdio {
		return ""
	}
	if c.PublicURL != "" {
		return c.PublicURL
	}
	host := strings.TrimSpace(c.BindHost)
	switch host {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%s%s", host, c.Port, c.ActiveEndpoint())
}

func normalizeTransport(raw string) transportMode {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "http", "streamable-http", "streamable_http", "streamablehttp":
		return transportModeStreamableHTTP
	case "sse":
		return transportModeSSE
	case "stdio":
		return transportModeStdio
	default:
		return transportModeStreamableHTTP
	}
}

func normalizeEndpointPath(path string) string {
	trimmed := strings.Trim(strings.TrimSpace(path), "/")
	if trimmed == "" {
		return "/"
	}
	return "/" + trimmed
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}

func logTransportStartup(cfg transportConfig, authMode security.AuthMode, adminEnabled bool, tokenExchangeEnabled bool) {
	switch cfg.Mode {
	case transportModeSSE:
		log.Printf("MCP SSE server listening on %s%s (message endpoint %s)", cfg.ListenAddr(), cfg.SSEEndpoint, cfg.MessageEndpoint)
		log.Printf("MCP transport: sse")
	case transportModeStdio:
		log.Printf("MCP stdio server ready")
		log.Printf("MCP transport: stdio")
	default:
		log.Printf("MCP streamable-http server listening on %s%s", cfg.ListenAddr(), cfg.HTTPEndpoint)
		log.Printf("MCP transport: streamable-http (stateless=%t)", cfg.StatelessHTTP)
	}

	switch authMode {
	case security.AuthModeSharedToken:
		log.Printf("MCP auth mode: shared-token (%s)", security.APIAuthFallbackHeader)
	case security.AuthModeIssuedToken:
		log.Printf("MCP auth mode: issued-token")
		log.Printf("Token store: %s", security.TokenStorePathFromEnv())
		log.Printf("Credential source: token-bound dynamic credentials")
	case security.AuthModeRequestCredential:
		log.Printf("MCP auth mode: request-credential")
		if cfg.Mode == transportModeStdio {
			log.Printf("Credential source: stdio process env (%s / %s / %s)", security.RequestCredentialClientIDEnv, security.RequestCredentialClientKeyEnv, security.RequestCredentialClientTokenEnv)
		} else {
			log.Printf("Credential source: %s / %s / %s", security.RequestSecretIDHeader, security.RequestSecretKeyHeader, security.RequestSessionTokenHeader)
		}
		if security.RequestCredentialValidateIdentityFromEnv() {
			log.Printf("Request credential identity validation: enabled via STS GetCallerIdentity")
		} else {
			log.Printf("warning: request credential identity validation is disabled")
		}
	default:
		log.Printf("MCP auth mode: none")
	}

	if cfg.UsesHTTP() {
		if adminEnabled {
			log.Printf("Admin token API enabled on /admin/tokens")
		} else if authMode == security.AuthModeIssuedToken {
			log.Printf("warning: issued-token mode is enabled but %s is empty; admin token API is disabled", security.AdminAPITokenEnv)
		}
		if tokenExchangeEnabled {
			log.Printf("TencentCloud token exchange enabled on /auth/token-exchange/tencentcloud")
		}
		log.Printf("Health check endpoints enabled on /healthz and /readyz")
		return
	}

	if security.AdminAPITokenFromEnv() != "" {
		log.Printf("warning: admin token API is not available when MCP_TRANSPORT=stdio")
	}
	if getEnvBool(security.TokenExchangeEnabledEnv, false) {
		log.Printf("warning: token exchange HTTP endpoints are not available when MCP_TRANSPORT=stdio")
	}
}

func buildStdioContextFunc(authMode security.AuthMode) (server.StdioContextFunc, error) {
	switch authMode {
	case security.AuthModeIssuedToken:
		return nil, fmt.Errorf("stdio transport does not support issued-token auth on the current main branch")
	case security.AuthModeRequestCredential:
		credential, err := requestCredentialFromEnv()
		if err != nil {
			return nil, err
		}
		principal, err := buildStdioPrincipal(credential)
		if err != nil {
			return nil, err
		}
		return func(ctx context.Context) context.Context {
			ctx = security.WithRequestCredential(ctx, credential)
			if principal != nil {
				ctx = security.WithPrincipal(ctx, principal)
			}
			return ctx
		}, nil
	default:
		return func(ctx context.Context) context.Context { return ctx }, nil
	}
}

func requestCredentialFromEnv() (*security.RequestCredential, error) {
	headers := security.RequestCredentialHeadersFromEnv()
	secretID := strings.TrimSpace(headers[security.RequestSecretIDHeader])
	secretKey := strings.TrimSpace(headers[security.RequestSecretKeyHeader])
	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("stdio + request-credential requires %s/%s, or fallback envs %s/%s", security.RequestCredentialClientIDEnv, security.RequestCredentialClientKeyEnv, "MCP_SECRET_ID", "MCP_SECRET_KEY")
	}
	return &security.RequestCredential{
		SecretID:     secretID,
		SecretKey:    secretKey,
		SessionToken: strings.TrimSpace(headers[security.RequestSessionTokenHeader]),
	}, nil
}

func buildStdioPrincipal(credential *security.RequestCredential) (*security.Principal, error) {
	if credential == nil {
		return nil, nil
	}
	principal := &security.Principal{
		TokenID:        "stdio-request-credential",
		SubjectType:    "request-credential",
		SubjectID:      "request:" + security.MaskSecretID(credential.SecretID),
		TenantID:       "request-credential",
		DisplayName:    security.MaskSecretID(credential.SecretID),
		Scopes:         append([]string(nil), security.RequestCredentialScopesFromEnv()...),
		AllowedRegions: append([]string(nil), security.RequestCredentialAllowedRegionsFromEnv()...),
	}
	if !security.RequestCredentialValidateIdentityFromEnv() {
		return principal, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = security.WithRequestCredential(ctx, credential)
	provider := security.NewRequestHeaderCredentialProvider()
	tcCredential, err := provider.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve request credential for stdio failed: %w", err)
	}
	identity, err := security.LookupTencentCloudIdentity(ctx, tcCredential, getEnv(security.STSRegionEnv, defaultSTSRegion))
	if err != nil {
		return nil, fmt.Errorf("validate stdio request credential identity failed: %w", err)
	}
	subjectType, subjectID, tenantID := identity.ToPrincipalIdentity()
	principal.SubjectType = subjectType
	principal.SubjectID = subjectID
	principal.TenantID = tenantID
	if displayName := identity.DisplayName(); displayName != "" {
		principal.DisplayName = displayName
	}
	return principal, nil
}

func stdioChildEnv() []string {
	env := append([]string(nil), os.Environ()...)
	if !envHasKey(env, transportEnv) {
		env = append(env, transportEnv+"="+string(transportModeStdio))
	}
	return env
}

func envHasKey(env []string, key string) bool {
	prefix := key + "="
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			return true
		}
	}
	return false
}

func buildTencentCloudCredentialFromContext(ctx context.Context) (*common.Credential, error) {
	provider := security.NewRequestHeaderCredentialProvider()
	return provider.Resolve(ctx)
}
