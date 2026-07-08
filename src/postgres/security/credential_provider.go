package security

import (
	"context"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

type CredentialProvider interface {
	Resolve(ctx context.Context) (*common.Credential, error)
}
