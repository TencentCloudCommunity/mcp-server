package security

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type CredentialBinding struct {
	TokenID        string
	SubjectID      string
	Provider       string
	CredentialKind string
	EncryptedBody  string
	ExpiresAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var ErrCredentialBindingNotFound = fmt.Errorf("credential binding not found")

type CredentialStore interface {
	PutCredentialBinding(ctx context.Context, binding *CredentialBinding) error
	GetCredentialBinding(ctx context.Context, tokenID string) (*CredentialBinding, error)
	DeleteCredentialBinding(ctx context.Context, tokenID string) error
}

type storedTencentCloudCredential struct {
	SecretID      string     `json:"secret_id"`
	SecretKey     string     `json:"secret_key"`
	SessionToken  string     `json:"session_token,omitempty"`
	CredentialTag string     `json:"credential_tag,omitempty"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
}

func (c storedTencentCloudCredential) ToCredential() (*common.Credential, error) {
	secretID := strings.TrimSpace(c.SecretID)
	secretKey := strings.TrimSpace(c.SecretKey)
	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("stored credential is incomplete")
	}
	if token := strings.TrimSpace(c.SessionToken); token != "" {
		return common.NewTokenCredential(secretID, secretKey, token), nil
	}
	return common.NewCredential(secretID, secretKey), nil
}

func encryptTencentCloudCredential(cipher *CredentialCipher, payload storedTencentCloudCredential) (string, error) {
	if cipher == nil {
		return "", fmt.Errorf("credential cipher is not initialized")
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal credential payload: %w", err)
	}
	return cipher.EncryptString(string(raw))
}

func decryptTencentCloudCredential(cipher *CredentialCipher, encrypted string) (*storedTencentCloudCredential, error) {
	if cipher == nil {
		return nil, fmt.Errorf("credential cipher is not initialized")
	}
	plaintext, err := cipher.DecryptString(encrypted)
	if err != nil {
		return nil, err
	}
	var payload storedTencentCloudCredential
	if err := json.Unmarshal([]byte(plaintext), &payload); err != nil {
		return nil, fmt.Errorf("unmarshal credential payload: %w", err)
	}
	return &payload, nil
}
