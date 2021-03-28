package pkg

import (
	"amazing_talker/pkg/model"
	"context"
)

// EamilService service介面層
type EamilService interface {
	SendVerifyEmail(ctx context.Context, account model.IdentityAccount) error
}
