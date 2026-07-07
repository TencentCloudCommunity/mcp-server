package security

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)

const (
	TokenExchangeModeEnv                   = "MCP_TOKEN_EXCHANGE_MODE"
	legacyTokenExchangeModeEnv             = "TENCENTCLOUD_TOKEN_EXCHANGE_MODE"
	TokenExchangeEnabledEnv                = "TOKEN_EXCHANGE_ENABLED"
	TokenExchangeAllowedScopesEnv          = "TOKEN_EXCHANGE_ALLOWED_SCOPES"
	STSRegionEnv                           = "MCP_STS_REGION"
	legacySTSRegionEnv                     = "TENCENTCLOUD_STS_REGION"
	RoleArnEnv                             = "MCP_ROLE_ARN"
	legacyRoleArnEnv                       = "TENCENTCLOUD_ROLE_ARN"
	RoleSessionPrefixEnv                   = "MCP_ROLE_SESSION_PREFIX"
	legacyRoleSessionPrefixEnv             = "TENCENTCLOUD_ROLE_SESSION_PREFIX"
	RoleDurationEnv                        = "MCP_ROLE_DURATION_SECONDS"
	legacyRoleDurationEnv                  = "TENCENTCLOUD_ROLE_DURATION_SECONDS"
	defaultSTSRegion                       = "ap-guangzhou"
)

type TencentCloudTokenExchangeMode string

const (
	TokenExchangeModeSourceCredential TencentCloudTokenExchangeMode = "source-credential"
	TokenExchangeModeAssumeRole       TencentCloudTokenExchangeMode = "assume-role"
)

type TencentCloudTokenExchangeConfig struct {
	Enabled           bool
	Mode              TencentCloudTokenExchangeMode
	STSRegion         string
	AllowedScopes     []string
	RoleArn           string
	RoleSessionPrefix string
	RoleDuration      time.Duration
}

type TencentCloudIdentity struct {
	Arn         string
	AccountID   string
	UserID      string
	PrincipalID string
	Type        string
	RequestID   string
}

type TencentCloudTokenExchangeInput struct {
	SecretID        string
	SecretKey       string
	SessionToken    string
	DisplayName     string
	AllowedRegions  []string
	RequestedScopes []string
	ExpiresIn       time.Duration
	Description     string
}

type TencentCloudTokenExchangeResult struct {
	Issued              *IssuedToken
	Identity            TencentCloudIdentity
	CredentialKind      string
	CredentialExpiresAt *time.Time
}

type TencentCloudTokenExchangeService struct {
	issuer *TokenIssuer
	store  CredentialStore
	cipher *CredentialCipher
	cfg    TencentCloudTokenExchangeConfig
}

func TencentCloudTokenExchangeConfigFromEnv(authMode AuthMode) TencentCloudTokenExchangeConfig {
	enabled := authMode == AuthModeIssuedToken
	if raw := strings.TrimSpace(os.Getenv(TokenExchangeEnabledEnv)); raw != "" {
		enabled = strings.EqualFold(raw, "true") || raw == "1" || strings.EqualFold(raw, "yes")
	}
	cfg := TencentCloudTokenExchangeConfig{
		Enabled:           enabled,
		Mode:              TokenExchangeModeSourceCredential,
		STSRegion:         defaultString(firstNonEmptyEnv(STSRegionEnv, legacySTSRegionEnv), defaultSTSRegion),
		AllowedScopes:     cleanUniqueStrings(strings.Split(defaultString(strings.TrimSpace(os.Getenv(TokenExchangeAllowedScopesEnv)), "pg.read"), ",")),
		RoleArn:           firstNonEmptyEnv(RoleArnEnv, legacyRoleArnEnv),
		RoleSessionPrefix: defaultString(firstNonEmptyEnv(RoleSessionPrefixEnv, legacyRoleSessionPrefixEnv), "pg-mcp"),
		RoleDuration:      secondsEnvAliases(3600, RoleDurationEnv, legacyRoleDurationEnv),
	}
	if rawMode := firstNonEmptyEnv(TokenExchangeModeEnv, legacyTokenExchangeModeEnv); rawMode != "" {
		cfg.Mode = TencentCloudTokenExchangeMode(rawMode)
	}
	if len(cfg.AllowedScopes) == 0 {
		cfg.AllowedScopes = []string{"pg.read"}
	}
	return cfg
}

func NewTencentCloudTokenExchangeService(issuer *TokenIssuer, store CredentialStore, cipher *CredentialCipher, cfg TencentCloudTokenExchangeConfig) (*TencentCloudTokenExchangeService, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if issuer == nil {
		return nil, fmt.Errorf("token issuer is not initialized")
	}
	if store == nil {
		return nil, fmt.Errorf("credential store is not initialized")
	}
	if cipher == nil {
		return nil, fmt.Errorf("credential cipher is not initialized")
	}
	if cfg.STSRegion == "" {
		cfg.STSRegion = defaultSTSRegion
	}
	switch cfg.Mode {
	case TokenExchangeModeSourceCredential:
		// dev-compatible fallback; no extra validation.
	case TokenExchangeModeAssumeRole:
		if cfg.RoleArn == "" {
			return nil, fmt.Errorf("missing %s for assume-role token exchange", RoleArnEnv)
		}
	default:
		return nil, fmt.Errorf("unsupported %s=%q", TokenExchangeModeEnv, cfg.Mode)
	}
	if cfg.RoleDuration <= 0 {
		cfg.RoleDuration = time.Hour
	}
	if cfg.RoleDuration > 12*time.Hour {
		cfg.RoleDuration = 12 * time.Hour
	}
	return &TencentCloudTokenExchangeService{issuer: issuer, store: store, cipher: cipher, cfg: cfg}, nil
}

func (s *TencentCloudTokenExchangeService) Exchange(ctx context.Context, input TencentCloudTokenExchangeInput) (*TencentCloudTokenExchangeResult, error) {
	if s == nil {
		return nil, fmt.Errorf("token exchange service is not initialized")
	}
	sourceCredential, err := buildTencentCloudCredential(input.SecretID, input.SecretKey, input.SessionToken)
	if err != nil {
		return nil, err
	}
	identity, err := s.getCallerIdentity(ctx, sourceCredential)
	if err != nil {
		return nil, err
	}

	runtimeCredential := sourceCredential
	credentialKind := string(TokenExchangeModeSourceCredential)
	var credentialExpiresAt *time.Time
	if s.cfg.Mode == TokenExchangeModeAssumeRole {
		runtimeCredential, credentialExpiresAt, err = s.assumeRole(ctx, sourceCredential, identity)
		if err != nil {
			return nil, err
		}
		credentialKind = string(TokenExchangeModeAssumeRole)
	}

	subjectType, subjectID, tenantID := identity.ToPrincipalIdentity()
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" {
		displayName = identity.DisplayName()
	}
	scopes, err := s.normalizeScopes(input.RequestedScopes)
	if err != nil {
		return nil, err
	}
	expiresIn := input.ExpiresIn
	if credentialExpiresAt != nil {
		remaining := time.Until(*credentialExpiresAt)
		if remaining <= 0 {
			return nil, fmt.Errorf("temporary credential returned by sts has already expired")
		}
		if expiresIn <= 0 || remaining < expiresIn {
			expiresIn = remaining
		}
	}
	issued, err := s.issuer.Issue(ctx, CreateTokenInput{
		SubjectType:    subjectType,
		SubjectID:      subjectID,
		TenantID:       tenantID,
		DisplayName:    displayName,
		Scopes:         scopes,
		AllowedRegions: cleanUniqueStrings(input.AllowedRegions),
		ExpiresIn:      expiresIn,
		Description:    strings.TrimSpace(input.Description),
		IssuedBy:       "self-service/tencentcloud",
	})
	if err != nil {
		return nil, err
	}
	secretID, secretKey, sessionToken := runtimeCredential.GetCredential()
	encryptedBody, err := encryptTencentCloudCredential(s.cipher, storedTencentCloudCredential{
		SecretID:      secretID,
		SecretKey:     secretKey,
		SessionToken:  sessionToken,
		CredentialTag: credentialKind,
		ExpiresAt:     credentialExpiresAt,
	})
	if err != nil {
		s.revokeIssuedTokenBestEffort(ctx, issued.Record.ID)
		return nil, err
	}
	if err := s.store.PutCredentialBinding(ctx, &CredentialBinding{
		TokenID:        issued.Record.ID,
		SubjectID:      issued.Record.SubjectID,
		Provider:       "tencentcloud",
		CredentialKind: credentialKind,
		EncryptedBody:  encryptedBody,
		ExpiresAt:      credentialExpiresAt,
		CreatedAt:      issued.Record.CreatedAt,
		UpdatedAt:      issued.Record.UpdatedAt,
	}); err != nil {
		s.revokeIssuedTokenBestEffort(ctx, issued.Record.ID)
		return nil, err
	}
	return &TencentCloudTokenExchangeResult{
		Issued:              issued,
		Identity:            identity,
		CredentialKind:      credentialKind,
		CredentialExpiresAt: credentialExpiresAt,
	}, nil
}

func (s *TencentCloudTokenExchangeService) normalizeScopes(requested []string) ([]string, error) {
	allowed := make(map[string]string, len(s.cfg.AllowedScopes))
	for _, scope := range s.cfg.AllowedScopes {
		allowed[strings.ToLower(scope)] = scope
	}
	requested = cleanUniqueStrings(requested)
	if len(requested) == 0 {
		return append([]string(nil), s.cfg.AllowedScopes...), nil
	}
	result := make([]string, 0, len(requested))
	for _, scope := range requested {
		if normalized, ok := allowed[strings.ToLower(scope)]; ok {
			result = append(result, normalized)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("requested scopes are not allowed for self-service token exchange")
	}
	return result, nil
}

func (s *TencentCloudTokenExchangeService) getCallerIdentity(ctx context.Context, credential *common.Credential) (TencentCloudIdentity, error) {
	cpf := profile.NewClientProfile()
	client, err := sts.NewClient(credential, s.cfg.STSRegion, cpf)
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

func (s *TencentCloudTokenExchangeService) assumeRole(ctx context.Context, credential *common.Credential, identity TencentCloudIdentity) (*common.Credential, *time.Time, error) {
	cpf := profile.NewClientProfile()
	client, err := sts.NewClient(credential, s.cfg.STSRegion, cpf)
	if err != nil {
		return nil, nil, fmt.Errorf("create sts client for assume role failed: %w", err)
	}
	request := sts.NewAssumeRoleRequest()
	request.RoleArn = common.StringPtr(s.cfg.RoleArn)
	request.RoleSessionName = common.StringPtr(s.buildRoleSessionName(identity))
	durationSeconds := uint64(s.cfg.RoleDuration / time.Second)
	request.DurationSeconds = &durationSeconds
	response, err := client.AssumeRoleWithContext(ctx, request)
	if err != nil {
		return nil, nil, fmt.Errorf("assume role failed: %w", err)
	}
	if response == nil || response.Response == nil || response.Response.Credentials == nil {
		return nil, nil, fmt.Errorf("assume role returned empty credentials")
	}
	creds := response.Response.Credentials
	tmpSecretID := pointerString(creds.TmpSecretId)
	tmpSecretKey := pointerString(creds.TmpSecretKey)
	token := pointerString(creds.Token)
	if tmpSecretID == "" || tmpSecretKey == "" || token == "" {
		return nil, nil, fmt.Errorf("assume role returned incomplete temporary credentials")
	}
	expiration := pointerString(response.Response.Expiration)
	var expiresAt *time.Time
	if expiration != "" {
		if parsed, err := time.Parse(time.RFC3339, expiration); err == nil {
			expiresAt = &parsed
		}
	}
	return common.NewTokenCredential(tmpSecretID, tmpSecretKey, token), expiresAt, nil
}

func (s *TencentCloudTokenExchangeService) buildRoleSessionName(identity TencentCloudIdentity) string {
	prefix := defaultString(strings.TrimSpace(s.cfg.RoleSessionPrefix), "pg-mcp")
	base := identity.UserID
	if base == "" {
		base = identity.PrincipalID
	}
	if base == "" {
		base = "caller"
	}
	safe := invalidRoleSessionChars.ReplaceAllString(strings.ToLower(base), "-")
	safe = strings.Trim(safe, "-")
	if safe == "" {
		safe = "caller"
	}
	sessionName := fmt.Sprintf("%s-%s", prefix, safe)
	if len(sessionName) > 64 {
		sessionName = sessionName[:64]
	}
	return sessionName
}

func (i TencentCloudIdentity) ToPrincipalIdentity() (subjectType, subjectID, tenantID string) {
	accountID := defaultString(strings.TrimSpace(i.AccountID), "unknown-account")
	userID := defaultString(strings.TrimSpace(i.UserID), defaultString(strings.TrimSpace(i.PrincipalID), "unknown-user"))
	identityType := strings.ToLower(defaultString(strings.TrimSpace(i.Type), "user"))
	return "tencentcloud-" + identityType, fmt.Sprintf("tencentcloud:%s:%s", accountID, userID), fmt.Sprintf("tencentcloud-account:%s", accountID)
}

func (i TencentCloudIdentity) DisplayName() string {
	if arn := strings.TrimSpace(i.Arn); arn != "" {
		return arn
	}
	if userID := strings.TrimSpace(i.UserID); userID != "" {
		return userID
	}
	return strings.TrimSpace(i.PrincipalID)
}

func buildTencentCloudCredential(secretID, secretKey, sessionToken string) (*common.Credential, error) {
	secretID = strings.TrimSpace(secretID)
	secretKey = strings.TrimSpace(secretKey)
	sessionToken = strings.TrimSpace(sessionToken)
	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("secret_id and secret_key are required")
	}
	if sessionToken != "" {
		return common.NewTokenCredential(secretID, secretKey, sessionToken), nil
	}
	return common.NewCredential(secretID, secretKey), nil
}

func (s *TencentCloudTokenExchangeService) revokeIssuedTokenBestEffort(ctx context.Context, tokenID string) {
	tokenStore, ok := s.store.(TokenStore)
	if !ok {
		return
	}
	_ = tokenStore.RevokeToken(ctx, tokenID, RevokeTokenInput{RevokedBy: "system", Reason: "token exchange binding failed"})
}

func pointerString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func firstNonEmptyEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func secondsEnvAliases(fallback int, keys ...string) time.Duration {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			return time.Duration(fallback) * time.Second
		}
		return time.Duration(parsed) * time.Second
	}
	return time.Duration(fallback) * time.Second
}

var invalidRoleSessionChars = regexp.MustCompile(`[^a-z0-9+=,.@_-]+`)
