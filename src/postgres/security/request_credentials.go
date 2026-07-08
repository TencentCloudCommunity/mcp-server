package security

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)

const (
	RequestSecretIDHeader           = "X-TencentCloud-Secret-Id"
	RequestSecretKeyHeader          = "X-TencentCloud-Secret-Key"
	RequestSessionTokenHeader       = "X-TencentCloud-Session-Token"
	RequestCredentialScopesEnv      = "MCP_REQUEST_CREDENTIAL_SCOPES"
	RequestCredentialRegionsEnv     = "MCP_REQUEST_ALLOWED_REGIONS"
	RequestCredentialValidateEnv    = "MCP_REQUEST_VALIDATE_IDENTITY"
	RequestCredentialClientIDEnv    = "MCP_REQUEST_SECRET_ID"
	RequestCredentialClientKeyEnv   = "MCP_REQUEST_SECRET_KEY"
	RequestCredentialClientTokenEnv = "MCP_REQUEST_SESSION_TOKEN"
)

type RequestCredential struct {
	SecretID     string
	SecretKey    string
	SessionToken string
}

type requestCredentialContextKey struct{}

type RequestHeaderCredentialProvider struct{}

func WithRequestCredential(ctx context.Context, credential *RequestCredential) context.Context {
	if ctx == nil || credential == nil {
		return ctx
	}
	return context.WithValue(ctx, requestCredentialContextKey{}, credential)
}

func RequestCredentialFromContext(ctx context.Context) (*RequestCredential, bool) {
	if ctx == nil {
		return nil, false
	}
	credential, ok := ctx.Value(requestCredentialContextKey{}).(*RequestCredential)
	return credential, ok && credential != nil
}

func NewRequestHeaderCredentialProvider() *RequestHeaderCredentialProvider {
	return &RequestHeaderCredentialProvider{}
}

func (p *RequestHeaderCredentialProvider) Resolve(ctx context.Context) (*common.Credential, error) {
	credential, ok := RequestCredentialFromContext(ctx)
	if !ok || credential == nil {
		return nil, fmt.Errorf("request credential is not available in context")
	}
	return buildTencentCloudCredential(credential.SecretID, credential.SecretKey, credential.SessionToken)
}

func ExtractRequestCredential(r *http.Request) *RequestCredential {
	if r == nil {
		return nil
	}
	secretID := strings.TrimSpace(r.Header.Get(RequestSecretIDHeader))
	secretKey := strings.TrimSpace(r.Header.Get(RequestSecretKeyHeader))
	if secretID == "" || secretKey == "" {
		return nil
	}
	return &RequestCredential{
		SecretID:     secretID,
		SecretKey:    secretKey,
		SessionToken: strings.TrimSpace(r.Header.Get(RequestSessionTokenHeader)),
	}
}

func RequestCredentialScopesFromEnv() []string {
	scopes := cleanUniqueStrings(strings.Split(strings.TrimSpace(os.Getenv(RequestCredentialScopesEnv)), ","))
	if len(scopes) == 0 {
		return []string{"pg.read", "pg.write"}
	}
	return scopes
}

func RequestCredentialAllowedRegionsFromEnv() []string {
	return cleanUniqueStrings(strings.Split(strings.TrimSpace(os.Getenv(RequestCredentialRegionsEnv)), ","))
}

func RequestCredentialValidateIdentityFromEnv() bool {
	raw := strings.TrimSpace(os.Getenv(RequestCredentialValidateEnv))
	if raw == "" {
		return true
	}
	switch strings.ToLower(raw) {
	case "0", "false", "no", "off":
		return false
	default:
		return true
	}
}

func RequestCredentialHeadersFromEnv() map[string]string {
	secretID := firstNonEmptyEnv(RequestCredentialClientIDEnv, "MCP_SECRET_ID", "TENCENTCLOUD_SECRET_ID")
	secretKey := firstNonEmptyEnv(RequestCredentialClientKeyEnv, "MCP_SECRET_KEY", "TENCENTCLOUD_SECRET_KEY")
	if strings.TrimSpace(secretID) == "" || strings.TrimSpace(secretKey) == "" {
		return nil
	}
	headers := map[string]string{
		RequestSecretIDHeader:  strings.TrimSpace(secretID),
		RequestSecretKeyHeader: strings.TrimSpace(secretKey),
	}
	if sessionToken := firstNonEmptyEnv(RequestCredentialClientTokenEnv, "TENCENTCLOUD_SESSION_TOKEN"); strings.TrimSpace(sessionToken) != "" {
		headers[RequestSessionTokenHeader] = strings.TrimSpace(sessionToken)
	}
	return headers
}

func RequestCredentialPlaceholderHeaders() map[string]string {
	return map[string]string{
		RequestSecretIDHeader:     "<TENCENTCLOUD_SECRET_ID>",
		RequestSecretKeyHeader:    "<TENCENTCLOUD_SECRET_KEY>",
		RequestSessionTokenHeader: "<TENCENTCLOUD_SESSION_TOKEN_OPTIONAL>",
	}
}

func LookupTencentCloudIdentity(ctx context.Context, credential *common.Credential, region string) (TencentCloudIdentity, error) {
	region = strings.TrimSpace(region)
	if region == "" {
		region = defaultSTSRegion
	}
	cpf := profile.NewClientProfile()
	client, err := sts.NewClient(credential, region, cpf)
	if err != nil {
		return TencentCloudIdentity{}, fmt.Errorf("create sts client failed: %w", err)
	}
	response, err := client.GetCallerIdentityWithContext(ctx, sts.NewGetCallerIdentityRequest())
	if err != nil {
		return TencentCloudIdentity{}, fmt.Errorf("get caller identity failed: %w", err)
	}
	if response == nil || response.Response == nil {
		return TencentCloudIdentity{}, fmt.Errorf("get caller identity returned empty response")
	}
	return TencentCloudIdentity{
		Arn:         pointerString(response.Response.Arn),
		AccountID:   pointerString(response.Response.AccountId),
		UserID:      pointerString(response.Response.UserId),
		PrincipalID: pointerString(response.Response.PrincipalId),
		Type:        pointerString(response.Response.Type),
		RequestID:   pointerString(response.Response.RequestId),
	}, nil
}

func MaskSecretID(secretID string) string {
	secretID = strings.TrimSpace(secretID)
	if len(secretID) <= 8 {
		if secretID == "" {
			return "anonymous"
		}
		return secretID
	}
	return secretID[:4] + "***" + secretID[len(secretID)-4:]
}
