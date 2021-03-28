package testsuite

import (
	"amazing_talker/configuration"
	"amazing_talker/internal"
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
	"gitlab.com/howmay/gopher/db"
	"gitlab.com/howmay/gopher/zlog"
	"go.uber.org/fx"
)

// Suite ...
type Suite struct {
	app         *fx.App
	pgContainer *dockertest.Resource
	pool        *dockertest.Pool
	t           testing.T
}

var suite Suite

// Initialize 初始化 suite
func Initialize(fxOption ...fx.Option) error {
	if os.Getenv("CONFIG_NAME") == "" {
		_ = os.Setenv("CONFIG_NAME", "app-test")
	}
	configuration, err := configuration.New()
	if err != nil {
		return err
	}

	suite.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Error().Msgf("Could not connect to docker: %s", err)
		return err
	}
	databaseName := "postgres"
	suite.pgContainer, err = suite.pool.Run("postgres", "10-alpine", []string{"POSTGRES_PASSWORD=admin", "POSTGRES_DB=" + databaseName})
	if err != nil {
		log.Error().Msgf("Could not start resource: %s", err)
		return err
	}
	pgPort, _ := strconv.Atoi(suite.pgContainer.GetPort("5432/tcp"))

	dbCfg := &db.Database{
		Host:     "localhost",
		Port:     pgPort,
		Password: "admin",
		Username: "postgres",
		Type:     db.Postgres,
		DBName:   databaseName,
		Debug:    true,
	}
	configuration.Database.Read = dbCfg
	configuration.Database.Write = dbCfg

	base := []fx.Option{
		fx.Supply(*configuration),
		fx.Provide(
			db.InitDatabases,
		),
		fx.Invoke(db.SetEncrypKey),
		fx.Invoke(zlog.InitV2),
		fx.Invoke(internal.MigrationDB),
	}

	base = append(base, fxOption...)

	app := fx.New(
		base...,
	)

	suite.app = app
	return app.Start(context.Background())
}

// Close 停止 container
func Close() {
	log.Info().Msg("close app")
	if err := suite.pool.Purge(suite.pgContainer); err != nil {
		log.Error().Msgf("Could not purge resource: %s", err)
	}
}
