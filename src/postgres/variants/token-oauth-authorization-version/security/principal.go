package security

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Principal struct {
	TokenID        string
	SubjectType    string
	SubjectID      string
	TenantID       string
	DisplayName    string
	Scopes         []string
	AllowedRegions []string
	ExpiresAt      time.Time
}

type principalContextKey struct{}

func WithPrincipal(ctx context.Context, principal *Principal) context.Context {
	if ctx == nil || principal == nil {
		return ctx
	}
	return context.WithValue(ctx, principalContextKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (*Principal, bool) {
	if ctx == nil {
		return nil, false
	}
	principal, ok := ctx.Value(principalContextKey{}).(*Principal)
	return principal, ok && principal != nil
}

func (p *Principal) HasScope(scope string) bool {
	scope = strings.TrimSpace(scope)
	if scope == "" || p == nil {
		return true
	}
	for _, candidate := range p.Scopes {
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || strings.EqualFold(candidate, scope) {
			return true
		}
	}
	return false
}

func (p *Principal) AllowsRegion(region string) bool {
	region = strings.TrimSpace(region)
	if p == nil || region == "" || len(p.AllowedRegions) == 0 {
		return true
	}
	for _, allowed := range p.AllowedRegions {
		if strings.EqualFold(strings.TrimSpace(allowed), region) {
			return true
		}
	}
	return false
}

func RequiredScopeForGuardLevel(level GuardLevel) string {
	if level == LevelNone {
		return "pg.read"
	}
	return "pg.write"
}

func AuthorizePrincipal(ctx context.Context, toolName, region string, level GuardLevel) error {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal == nil {
		return nil
	}
	if !principal.ExpiresAt.IsZero() && time.Now().After(principal.ExpiresAt) {
		return fmt.Errorf("token for subject '%s' has expired", principal.SubjectID)
	}
	requiredScope := RequiredScopeForGuardLevel(level)
	if !principal.HasScope(requiredScope) {
		return fmt.Errorf("subject '%s' lacks scope '%s' for tool '%s'", principal.SubjectID, requiredScope, toolName)
	}
	if !principal.AllowsRegion(region) {
		return fmt.Errorf("region '%s' is outside allowed scope for subject '%s'", region, principal.SubjectID)
	}
	return nil
}
