package security

import (
	"context"
	"errors"
	"time"
)

var ErrTokenNotFound = errors.New("token not found")

type TokenStatus string

const (
	TokenStatusActive  TokenStatus = "active"
	TokenStatusRevoked TokenStatus = "revoked"
	TokenStatusExpired TokenStatus = "expired"
)

type TokenRecord struct {
	ID             string
	TokenHash      string
	TokenPrefix    string
	SubjectType    string
	SubjectID      string
	TenantID       string
	DisplayName    string
	Scopes         []string
	AllowedRegions []string
	Status         TokenStatus
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastUsedAt     *time.Time
	IssuedBy       string
	RevokedBy      string
	RevokedAt      *time.Time
	RevokeReason   string
	Description    string
}

type TokenFilter struct {
	TenantID  string
	SubjectID string
	Status    TokenStatus
	Scope     string
}

type CreateTokenInput struct {
	SubjectType    string
	SubjectID      string
	TenantID       string
	DisplayName    string
	Scopes         []string
	AllowedRegions []string
	ExpiresIn      time.Duration
	Description    string
	IssuedBy       string
}

type RevokeTokenInput struct {
	RevokedBy string
	Reason    string
}

type IssuedToken struct {
	Token  string
	Record TokenRecord
}

type TokenStore interface {
	CreateToken(ctx context.Context, record *TokenRecord) error
	GetTokenByID(ctx context.Context, id string) (*TokenRecord, error)
	GetTokenByHash(ctx context.Context, tokenHash string) (*TokenRecord, error)
	ListTokens(ctx context.Context, filter TokenFilter) ([]TokenRecord, error)
	RevokeToken(ctx context.Context, id string, input RevokeTokenInput) error
	TouchToken(ctx context.Context, id string, usedAt time.Time) error
}

func (r TokenRecord) ToPrincipal() *Principal {
	return &Principal{
		TokenID:        r.ID,
		SubjectType:    r.SubjectType,
		SubjectID:      r.SubjectID,
		TenantID:       r.TenantID,
		DisplayName:    r.DisplayName,
		Scopes:         append([]string(nil), r.Scopes...),
		AllowedRegions: append([]string(nil), r.AllowedRegions...),
		ExpiresAt:      r.ExpiresAt,
	}
}
