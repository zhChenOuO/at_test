package service

import (
	"amazing_talker/internal/bundle"
	"amazing_talker/internal/igmail"
	"amazing_talker/internal/testsuite"
	"amazing_talker/pkg"
	"amazing_talker/pkg/repository"
	"os"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"google.golang.org/api/gmail/v1"
)

type Suite struct {
	repo     pkg.IRepository
	gmailSvc *gmail.Service
	bundle   *i18n.Bundle
}

var suite Suite

func TestMain(m *testing.M) {
	testsuite.Initialize(
		fx.Provide(
			repository.NewRepository,
			igmail.NewGmailService,
			bundle.NewBundle,
			NewService,
		),
		fx.Populate(&suite.repo, &suite.gmailSvc, &suite.bundle),
	)
	code := m.Run()
	testsuite.Close()

	os.Exit(code)
}
