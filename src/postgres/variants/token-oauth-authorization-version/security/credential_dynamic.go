package security

import (
	"context"
	"fmt"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type TokenBoundCredentialProvider struct {
	store  CredentialStore
	cipher *CredentialCipher
}

func NewTokenBoundCredentialProvider(store CredentialStore, cipher *CredentialCipher) (*TokenBoundCredentialProvider, error) {
	if store == nil {
		return nil, fmt.Errorf("credential store is not initialized")
	}
	if cipher == nil {
		return nil, fmt.Errorf("credential cipher is not initialized")
	}
	return &TokenBoundCredentialProvider{store: store, cipher: cipher}, nil
}

func (p *TokenBoundCredentialProvider) Resolve(ctx context.Context) (*common.Credential, error) {
	principal, ok := PrincipalFromContext(ctx)
	if !ok || principal == nil {
		return nil, fmt.Errorf("principal is missing from request context")
	}
	binding, err := p.store.GetCredentialBinding(ctx, principal.TokenID)
	if err != nil {
		return nil, fmt.Errorf("load credential binding: %w", err)
	}
	if binding.SubjectID != "" && binding.SubjectID != principal.SubjectID {
		return nil, fmt.Errorf("credential binding subject mismatch for token %s", principal.TokenID)
	}
	if binding.ExpiresAt != nil && !binding.ExpiresAt.IsZero() && time.Now().After(*binding.ExpiresAt) {
		return nil, fmt.Errorf("bound cloud credential for subject '%s' has expired", principal.SubjectID)
	}
	payload, err := decryptTencentCloudCredential(p.cipher, binding.EncryptedBody)
	if err != nil {
		return nil, fmt.Errorf("decrypt bound credential failed: %w", err)
	}
	if payload.ExpiresAt != nil && !payload.ExpiresAt.IsZero() && time.Now().After(*payload.ExpiresAt) {
		return nil, fmt.Errorf("bound cloud credential for subject '%s' has expired", principal.SubjectID)
	}
	return payload.ToCredential()
}
