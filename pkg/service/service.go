package service

import (
	pkg "amazing_talker/pkg"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/api/gmail/v1"
)

type service struct {
	repo     pkg.IRepository
	gmailSvc *gmail.Service
	bundle   *i18n.Bundle
}

// NewService 依賴注入
func NewService(repo pkg.IRepository, gmailSvc *gmail.Service, bundle *i18n.Bundle) pkg.IService {
	return &service{repo: repo, gmailSvc: gmailSvc, bundle: bundle}
}
