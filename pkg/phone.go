package pkg

import (
	"amazing_talker/pkg/model"
	"context"
)

// PhoneService service介面層
type PhoneService interface {
	SendVerifyPhone(ctx context.Context, account model.IdentityAccount) error
}
