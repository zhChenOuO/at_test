package internal

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/cenk/backoff"
	"github.com/pressly/goose"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.com/howmay/gopher/db"
)

// MigrationDB 初始化
func MigrationDB(dbCfg *db.Config) error {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		dbCfg.Write.Username, dbCfg.Write.Password, dbCfg.Write.Host+":"+strconv.Itoa(dbCfg.Write.Port), dbCfg.Write.DBName)

	err := goose.SetDialect("postgres")
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var db *sql.DB
	err = backoff.Retry(func() error {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		err = db.Ping()
		if err != nil {
			log.Error().Msgf("main: %s ping error: %v", dbCfg.Write.Type, err)
			return err
		}
		return nil
	}, bo)

	if err != nil {
		log.Error().Msgf("main: %s connect err: %s", dbCfg.Write.Type, err.Error())
		return err
	}

	if err := goose.Run("up", db, viper.GetString("PROJ_DIR")+"/deploy/database"); err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}

// MigrationTestData 初始化測試資料
func MigrationTestData(dbCfg *db.Config) error {
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s/%s",
		dbCfg.Write.Username, dbCfg.Write.Password, dbCfg.Write.Host+":"+strconv.Itoa(dbCfg.Write.Port), dbCfg.Write.DBName)

	err := goose.SetDialect("postgres")
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var db *sql.DB
	err = backoff.Retry(func() error {
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		err = db.Ping()
		if err != nil {
			log.Error().Msgf("main: %s ping error: %v", "mysql", err)
			return err
		}
		return nil
	}, bo)

	if err != nil {
		log.Error().Msgf("main: mysql connect err: %s", err.Error())
		return err
	}

	if err := goose.Run("up", db, viper.GetString("PROJ_DIR")+"/test/data"); err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}
