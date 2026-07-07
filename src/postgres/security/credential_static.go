package security

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type StaticCredentialProvider struct {
	credential     *common.Credential
	usingLegacyEnv bool
}

func NewStaticCredentialProviderFromEnv() (*StaticCredentialProvider, error) {
	secretID, secretKey, usingLegacyEnv := resolveCloudCredentials()
	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("missing credentials: set MCP_SECRET_ID/MCP_SECRET_KEY")
	}
	return &StaticCredentialProvider{
		credential:     common.NewCredential(secretID, secretKey),
		usingLegacyEnv: usingLegacyEnv,
	}, nil
}

func (p *StaticCredentialProvider) Resolve(ctx context.Context) (*common.Credential, error) {
	_ = ctx
	if p == nil || p.credential == nil {
		return nil, fmt.Errorf("credential provider is not initialized")
	}
	return p.credential, nil
}

func (p *StaticCredentialProvider) UsingLegacyEnv() bool {
	return p != nil && p.usingLegacyEnv
}

func resolveCloudCredentials() (secretID, secretKey string, usingLegacyEnv bool) {
	secretID = strings.TrimSpace(os.Getenv("MCP_SECRET_ID"))
	secretKey = strings.TrimSpace(os.Getenv("MCP_SECRET_KEY"))
	if secretID == "" {
		if legacy := strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_ID")); legacy != "" {
			secretID = legacy
			usingLegacyEnv = true
		}
	}
	if secretKey == "" {
		if legacy := strings.TrimSpace(os.Getenv("TENCENTCLOUD_SECRET_KEY")); legacy != "" {
			secretKey = legacy
			usingLegacyEnv = true
		}
	}
	return secretID, secretKey, usingLegacyEnv
}
