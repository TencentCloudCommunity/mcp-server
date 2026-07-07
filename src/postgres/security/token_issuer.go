package security

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	TokenStoreEnv          = "TOKEN_STORE"
	TokenStorePathEnv      = "TOKEN_STORE_PATH"
	TokenHashPepperEnv     = "TOKEN_HASH_PEPPER"
	TokenDefaultTTLEnv     = "TOKEN_DEFAULT_TTL_SECONDS"
	TokenMaxTTLEnv         = "TOKEN_MAX_TTL_SECONDS"
	TokenDefaultScopesEnv  = "TOKEN_DEFAULT_SCOPES"
	AdminAPITokenEnv       = "ADMIN_API_TOKEN"
	AdminActorHeader       = "X-Admin-Actor"
	defaultTokenStorePath  = "data/tokens.db"
	defaultTokenPrefix     = "mcp_ptk_"
	defaultTokenTTL        = 24 * time.Hour
	defaultTokenMaxTTL     = 30 * 24 * time.Hour
	defaultTokenEntropyLen = 32
)

type TokenIssuerConfig struct {
	Pepper        string
	DefaultTTL    time.Duration
	MaxTTL        time.Duration
	DefaultScopes []string
	Clock         func() time.Time
}

type TokenIssuer struct {
	store TokenStore
	cfg   TokenIssuerConfig
}

func NewTokenIssuer(store TokenStore, cfg TokenIssuerConfig) *TokenIssuer {
	if cfg.DefaultTTL <= 0 {
		cfg.DefaultTTL = defaultTokenTTL
	}
	if cfg.MaxTTL <= 0 {
		cfg.MaxTTL = defaultTokenMaxTTL
	}
	if len(cfg.DefaultScopes) == 0 {
		cfg.DefaultScopes = []string{"pg.read"}
	}
	if cfg.Clock == nil {
		cfg.Clock = func() time.Time { return time.Now().UTC() }
	}
	return &TokenIssuer{store: store, cfg: cfg}
}

func TokenIssuerConfigFromEnv() TokenIssuerConfig {
	return TokenIssuerConfig{
		Pepper:        strings.TrimSpace(os.Getenv(TokenHashPepperEnv)),
		DefaultTTL:    secondsEnv(TokenDefaultTTLEnv, int(defaultTokenTTL/time.Second)),
		MaxTTL:        secondsEnv(TokenMaxTTLEnv, int(defaultTokenMaxTTL/time.Second)),
		DefaultScopes: cleanUniqueStrings(strings.Split(strings.TrimSpace(os.Getenv(TokenDefaultScopesEnv)), ",")),
	}
}

func AdminAPITokenFromEnv() string {
	return strings.TrimSpace(os.Getenv(AdminAPITokenEnv))
}

func TokenStorePathFromEnv() string {
	if path := strings.TrimSpace(os.Getenv(TokenStorePathEnv)); path != "" {
		return path
	}
	return defaultTokenStorePath
}

func (i *TokenIssuer) Issue(ctx context.Context, input CreateTokenInput) (*IssuedToken, error) {
	if i == nil || i.store == nil {
		return nil, fmt.Errorf("token issuer is not initialized")
	}
	subjectID := strings.TrimSpace(input.SubjectID)
	if subjectID == "" {
		return nil, fmt.Errorf("subject_id is required")
	}
	now := i.cfg.Clock()
	ttl := input.ExpiresIn
	if ttl <= 0 {
		ttl = i.cfg.DefaultTTL
	}
	if ttl <= 0 {
		ttl = defaultTokenTTL
	}
	if i.cfg.MaxTTL > 0 && ttl > i.cfg.MaxTTL {
		ttl = i.cfg.MaxTTL
	}
	scopes := cleanUniqueStrings(input.Scopes)
	if len(scopes) == 0 {
		scopes = append([]string(nil), i.cfg.DefaultScopes...)
	}
	if len(scopes) == 0 {
		scopes = []string{"pg.read"}
	}
	token, err := generateOpaqueToken()
	if err != nil {
		return nil, err
	}
	record := &TokenRecord{
		ID:             "tok_" + uuid.NewString(),
		TokenHash:      HashToken(token, i.cfg.Pepper),
		TokenPrefix:    tokenPrefix(token),
		SubjectType:    defaultString(strings.TrimSpace(input.SubjectType), "user"),
		SubjectID:      subjectID,
		TenantID:       defaultString(strings.TrimSpace(input.TenantID), "default"),
		DisplayName:    strings.TrimSpace(input.DisplayName),
		Scopes:         scopes,
		AllowedRegions: cleanUniqueStrings(input.AllowedRegions),
		Status:         TokenStatusActive,
		ExpiresAt:      now.Add(ttl),
		CreatedAt:      now,
		UpdatedAt:      now,
		IssuedBy:       strings.TrimSpace(input.IssuedBy),
		Description:    strings.TrimSpace(input.Description),
	}
	if err := i.store.CreateToken(ctx, record); err != nil {
		return nil, err
	}
	return &IssuedToken{Token: token, Record: *record}, nil
}

func HashToken(token, pepper string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(pepper) + strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

func cleanUniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, value)
	}
	return result
}

func secondsEnv(key string, fallback int) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return time.Duration(fallback) * time.Second
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return time.Duration(fallback) * time.Second
	}
	return time.Duration(parsed) * time.Second
}

func generateOpaqueToken() (string, error) {
	buf := make([]byte, defaultTokenEntropyLen)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return defaultTokenPrefix + base64.RawURLEncoding.EncodeToString(buf), nil
}

func tokenPrefix(token string) string {
	if len(token) <= 16 {
		return token
	}
	return token[:16]
}

func defaultString(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}
