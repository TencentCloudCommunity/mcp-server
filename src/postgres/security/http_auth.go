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
	AuthModeNone              AuthMode = "none"
	AuthModeSharedToken       AuthMode = "shared-token"
	AuthModeIssuedToken       AuthMode = "issued-token"
	AuthModeRequestCredential AuthMode = "request-credential"
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

type requestCredentialAuthenticator struct {
	stsRegion        string
	scopes           []string
	allowedRegions   []string
	validateIdentity bool
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
	case AuthModeNone, AuthModeSharedToken, AuthModeIssuedToken, AuthModeRequestCredential:
		return AuthMode(mode)
	}
	if APIAuthTokenFromEnv() != "" {
		return AuthModeSharedToken
	}
	if headers := RequestCredentialHeadersFromEnv(); len(headers) > 0 {
		return AuthModeRequestCredential
	}
	return AuthModeRequestCredential
}

func MCPClientHeadersFromEnv() map[string]string {
	mode := MCPAuthModeFromEnv()
	if mode == AuthModeRequestCredential {
		if headers := RequestCredentialHeadersFromEnv(); len(headers) > 0 {
			return headers
		}
	}
	if token := DataPlaneTokenFromEnv(); token != "" {
		return map[string]string{
			"Authorization": "Bearer " + token,
		}
	}
	if headers := RequestCredentialHeadersFromEnv(); len(headers) > 0 {
		return headers
	}
	return nil
}

func MCPClientOptionsFromEnv() []transport.ClientOption {
	if headers := MCPClientHeadersFromEnv(); len(headers) > 0 {
		return []transport.ClientOption{transport.WithHeaders(headers)}
	}
	return nil
}

func MCPStreamableHTTPClientOptionsFromEnv() []transport.StreamableHTTPCOption {
	if headers := MCPClientHeadersFromEnv(); len(headers) > 0 {
		return []transport.StreamableHTTPCOption{transport.WithHTTPHeaders(headers)}
	}
	return nil
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
	case AuthModeRequestCredential:
		return &requestCredentialAuthenticator{
			stsRegion:        defaultString(firstNonEmptyEnv(STSRegionEnv, legacySTSRegionEnv), defaultSTSRegion),
			scopes:           RequestCredentialScopesFromEnv(),
			allowedRegions:   RequestCredentialAllowedRegionsFromEnv(),
			validateIdentity: RequestCredentialValidateIdentityFromEnv(),
		}, nil
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
				if authenticator.Mode() == AuthModeSharedToken || authenticator.Mode() == AuthModeIssuedToken {
					w.Header().Set("WWW-Authenticate", `Bearer realm="mcp-server-postgres"`)
				}
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		if authenticator.Mode() == AuthModeRequestCredential {
			if credential := ExtractRequestCredential(r); credential != nil {
				ctx = WithRequestCredential(ctx, credential)
			}
		}
		if principal != nil {
			ctx = WithPrincipal(ctx, principal)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
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

func (a *requestCredentialAuthenticator) Authenticate(r *http.Request) (*Principal, error) {
	requestCredential := ExtractRequestCredential(r)
	if requestCredential == nil {
		return nil, ErrUnauthorized
	}
	credential, err := buildTencentCloudCredential(requestCredential.SecretID, requestCredential.SecretKey, requestCredential.SessionToken)
	if err != nil {
		return nil, ErrUnauthorized
	}
	principal := &Principal{
		TokenID:        "request-credential",
		SubjectType:    "request-credential",
		SubjectID:      "request:" + MaskSecretID(requestCredential.SecretID),
		TenantID:       "request-credential",
		DisplayName:    MaskSecretID(requestCredential.SecretID),
		Scopes:         append([]string(nil), a.scopes...),
		AllowedRegions: append([]string(nil), a.allowedRegions...),
	}
	if !a.validateIdentity {
		return principal, nil
	}
	identity, err := LookupTencentCloudIdentity(r.Context(), credential, a.stsRegion)
	if err != nil {
		return nil, ErrUnauthorized
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

func (a *requestCredentialAuthenticator) Mode() AuthMode {
	return AuthModeRequestCredential
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
