package security

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client/transport"
)

const (
	APIAuthTokenEnv       = "MCP_API_TOKEN"
	AccessTokenEnv        = "MCP_ACCESS_TOKEN"
	APIAuthFallbackHeader = "X-MCP-API-Token"
	MCPAuthModeEnv        = "MCP_AUTH_MODE"
	SharedTokenScopesEnv  = "MCP_SHARED_TOKEN_SCOPES"
)

type AuthMode string

const (
	AuthModeNone        AuthMode = "none"
	AuthModeSharedToken AuthMode = "shared-token"
	AuthModeIssuedToken AuthMode = "issued-token"
)

var ErrUnauthorized = errors.New("unauthorized")

type Authenticator interface {
	Authenticate(r *http.Request) (*Principal, error)
	Mode() AuthMode
}

type noopAuthenticator struct{}

type sharedTokenAuthenticator struct {
	expectedToken string
	scopes        []string
}

type issuedTokenAuthenticator struct {
	store  TokenStore
	pepper string
}

func APIAuthTokenFromEnv() string {
	return strings.TrimSpace(os.Getenv(APIAuthTokenEnv))
}

func DataPlaneTokenFromEnv() string {
	if token := strings.TrimSpace(os.Getenv(AccessTokenEnv)); token != "" {
		return token
	}
	return APIAuthTokenFromEnv()
}

func MCPAuthModeFromEnv() AuthMode {
	mode := strings.TrimSpace(os.Getenv(MCPAuthModeEnv))
	switch AuthMode(mode) {
	case AuthModeNone, AuthModeSharedToken, AuthModeIssuedToken:
		return AuthMode(mode)
	}
	if APIAuthTokenFromEnv() != "" {
		return AuthModeSharedToken
	}
	return AuthModeNone
}

func MCPClientOptionsFromEnv() []transport.ClientOption {
	token := DataPlaneTokenFromEnv()
	if token == "" {
		return nil
	}
	return []transport.ClientOption{
		transport.WithHeaders(map[string]string{
			"Authorization": "Bearer " + token,
		}),
	}
}

func SharedTokenScopesFromEnv() []string {
	scopes := cleanUniqueStrings(strings.Split(strings.TrimSpace(os.Getenv(SharedTokenScopesEnv)), ","))
	if len(scopes) == 0 {
		return []string{"pg.read", "pg.write"}
	}
	return scopes
}

func NewMCPAuthenticator(mode AuthMode, store TokenStore, pepper string) (Authenticator, error) {
	if mode == "" {
		mode = MCPAuthModeFromEnv()
	}
	switch mode {
	case AuthModeNone:
		return noopAuthenticator{}, nil
	case AuthModeSharedToken:
		expectedToken := APIAuthTokenFromEnv()
		if expectedToken == "" {
			return nil, fmt.Errorf("missing %s for shared-token auth", APIAuthTokenEnv)
		}
		return &sharedTokenAuthenticator{expectedToken: expectedToken, scopes: SharedTokenScopesFromEnv()}, nil
	case AuthModeIssuedToken:
		if store == nil {
			return nil, fmt.Errorf("token store is required for issued-token auth")
		}
		return &issuedTokenAuthenticator{store: store, pepper: pepper}, nil
	default:
		return nil, fmt.Errorf("unsupported auth mode %q", mode)
	}
}

func WrapMCPAuth(authenticator Authenticator, next http.Handler) http.Handler {
	if next == nil {
		return nil
	}
	if authenticator == nil || authenticator.Mode() == AuthModeNone {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, err := authenticator.Authenticate(r)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				w.Header().Set("WWW-Authenticate", `Bearer realm="mcp-server-postgres"`)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if principal != nil {
			r = r.WithContext(WithPrincipal(r.Context(), principal))
		}
		next.ServeHTTP(w, r)
	})
}

func WrapAdminAuth(next http.Handler) http.Handler {
	if next == nil {
		return nil
	}
	expectedToken := AdminAPITokenFromEnv()
	if expectedToken == "" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authorizedRequest(r, expectedToken) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="postgres-admin"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (noopAuthenticator) Authenticate(r *http.Request) (*Principal, error) {
	_ = r
	return nil, nil
}

func (noopAuthenticator) Mode() AuthMode {
	return AuthModeNone
}

func (a *sharedTokenAuthenticator) Authenticate(r *http.Request) (*Principal, error) {
	if !authorizedRequest(r, a.expectedToken) {
		return nil, ErrUnauthorized
	}
	return &Principal{
		TokenID:     "shared-token",
		SubjectType: "service",
		SubjectID:   "shared-token",
		TenantID:    "shared-token",
		DisplayName: "Shared MCP token",
		Scopes:      append([]string(nil), a.scopes...),
		ExpiresAt:   time.Time{},
	}, nil
}

func (a *sharedTokenAuthenticator) Mode() AuthMode {
	return AuthModeSharedToken
}

func (a *issuedTokenAuthenticator) Authenticate(r *http.Request) (*Principal, error) {
	token := extractRequestToken(r)
	if token == "" {
		return nil, ErrUnauthorized
	}
	record, err := a.store.GetTokenByHash(r.Context(), HashToken(token, a.pepper))
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}
	if record.Status != TokenStatusActive {
		return nil, ErrUnauthorized
	}
	if !record.ExpiresAt.IsZero() && time.Now().After(record.ExpiresAt) {
		return nil, ErrUnauthorized
	}
	_ = a.store.TouchToken(r.Context(), record.ID, time.Now().UTC())
	return record.ToPrincipal(), nil
}

func (a *issuedTokenAuthenticator) Mode() AuthMode {
	return AuthModeIssuedToken
}

func authorizedRequest(r *http.Request, expectedToken string) bool {
	if r == nil {
		return false
	}
	return secureTokenEqual(extractRequestToken(r), expectedToken)
}

func extractRequestToken(r *http.Request) string {
	if r == nil {
		return ""
	}
	if bearerToken := extractBearerToken(r.Header.Get("Authorization")); bearerToken != "" {
		return bearerToken
	}
	return strings.TrimSpace(r.Header.Get(APIAuthFallbackHeader))
}

func extractBearerToken(authHeader string) string {
	fields := strings.Fields(strings.TrimSpace(authHeader))
	if len(fields) != 2 {
		return ""
	}
	if !strings.EqualFold(fields[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(fields[1])
}

func secureTokenEqual(got, expected string) bool {
	got = strings.TrimSpace(got)
	expected = strings.TrimSpace(expected)
	if got == "" || expected == "" || len(got) != len(expected) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(expected)) == 1
}
