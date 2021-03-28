package cmd

import (
	"amazing_talker/configuration"
	"amazing_talker/internal/bundle"
	"amazing_talker/internal/igmail"
	"amazing_talker/pkg/delivery/scheduler"
	repository "amazing_talker/pkg/repository"
	service "amazing_talker/pkg/service"
	"context"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/rs/zerolog/log"
	cobra "github.com/spf13/cobra"
	"gitlab.com/howmay/gopher/db"
	"gitlab.com/howmay/gopher/zlog"
	fx "go.uber.org/fx"
)

// SchedulerCmd 是此程式的Service入口點
var SchedulerCmd = &cobra.Command{
	Run: schedulerRun,
	Use: "scheduler",
}

func schedulerRun(cmd *cobra.Command, args []string) {
	defer cmdRecover()

	rand.Seed(time.Now().UnixNano())

	config, err := configuration.New()
	if err != nil {
		os.Exit(0)
		return
	}

	zlog.InitV2(config.Log)

	s := &scheduler.Scheduler{}
	waitGroup := &sync.WaitGroup{}

	app := fx.New(
		fx.Supply(*config, waitGroup),
		fx.Provide(
			db.InitDatabases,
			repository.NewRepository,
			igmail.NewGmailService,
			bundle.NewBundle,
			service.NewService,
			scheduler.NewScheduler,
		),
		fx.Invoke(db.SetEncrypKey),
		fx.Populate(&s),
	)

	exitCode := 0
	if err := app.Start(context.Background()); err != nil {
		log.Error().Stack().Err(err).Msg("main: fx app Start failed")
		os.Exit(exitCode)
		return
	}
	ctx, cronJobCancel := context.WithCancel(context.Background())
	s.Start(ctx)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info().Msg("main: shutting down server...")
	cronJobCancel()

	waitSchedulerCh := make(chan struct{})
	go func(ch chan struct{}) {
		waitGroup.Wait()
		ch <- struct{}{}
	}(waitSchedulerCh)

	select {
	case <-time.After(8 * time.Minute):
		break
	case <-waitSchedulerCh:
		break
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Error().Stack().Err(err).Msg("main: fx app Stop failed")
		exitCode = 1
	}

	os.Exit(exitCode)
}
