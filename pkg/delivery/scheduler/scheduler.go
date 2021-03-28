package scheduler

import (
	"amazing_talker/configuration"
	"amazing_talker/pkg"
	"context"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"gitlab.com/howmay/gopher/middleware"
)

// NewScheduler ...
func NewScheduler(
	waitGroup *sync.WaitGroup,
	svc pkg.IService,
	cfg *configuration.App,
) (*Scheduler, error) {
	return &Scheduler{
		waitGroup: waitGroup,
		svc:       svc,
		cfg:       cfg,
	}, nil
}

// Scheduler ...
type Scheduler struct {
	waitGroup *sync.WaitGroup
	cfg       *configuration.App
	svc       pkg.IService
}

// Start ...
func (s *Scheduler) Start(ctx context.Context) error {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger), // or use cron.DefaultLogger
	))

	s.AddFunc(ctx, c, s.cfg.Scheduler.VerifyIdentityAccountEmailFreq, "verifyIdentityAccountEmail", s.verifyIdentityAccountEmail)
	s.AddFunc(ctx, c, s.cfg.Scheduler.VerifyIdentityAccountPhoneFreq, "verifyIdentityAccountPhone", s.verifyIdentityAccountPhone)

	c.Start()
	return nil
}

type schedulerCMD func(context.Context) error

// AddFunc ...
func (s *Scheduler) AddFunc(ctx context.Context, c *cron.Cron, spec string, endpointName string, cmd schedulerCMD) {
	// 用來控制scheduler不重複執行的channel,同時間只能有一個相同的scheduler執行
	ch := make(chan struct{}, 1)

	logger := log.With().Str("endpoint", endpointName).Logger()

	fn := func() {
		select {
		case <-ctx.Done():
			logger.Warn().Msgf("scheduler: %s scheduler isn't executed, due to context is canceled", endpointName)
			return
		default:
			break
		}

		select {
		case ch <- struct{}{}:
			// 上次scheduler已經結束,繼續往下跑,不return
		default:
			logger.Warn().Msgf("scheduler: %s scheduler is pending", endpointName) // 還有相同scheduler再處理,忽略掉這次執行,直接return
			return
		}

		startTime := time.Now()
		requestID := uuid.New().String()

		ctx = middleware.ScheduleInitAddFuncAfterMiddleware(ctx)

		logger := logger.With().Fields(map[string]interface{}{
			"start_time": startTime.String(),
			"request_id": requestID,
		}).Logger()

		funcCtx := logger.WithContext(context.Background())

		defer func(startTime time.Time) {
			endTime := time.Now()
			if err := recover(); err != nil {
				trace := make([]byte, 4096)
				runtime.Stack(trace, true)
				var msg string
				for i := 2; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					msg += fmt.Sprintf("%s:%d\n", file, line)
				}

				logger.Error().Stack().Err(err.(error)).Fields(
					map[string]interface{}{
						"stack_error":   string(trace),
						"end_time":      endTime.String(),
						"latency":       strconv.FormatInt(int64(endTime.Sub(startTime)), 10),
						"latency_human": endTime.Sub(startTime).String(),
					}).Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥\n%s scheduler error", err, msg, endpointName)
			}
			<-ch
		}(startTime)

		s.waitGroup.Add(1)
		defer s.waitGroup.Done()
		err := cmd(funcCtx)
		endTime := time.Now()

		logFields := map[string]interface{}{
			"end_time":      endTime.String(),
			"latency":       strconv.FormatInt(int64(endTime.Sub(startTime)), 10),
			"latency_human": endTime.Sub(startTime).String(),
		}
		if err != nil {
			logger.Error().Stack().Err(err).Fields(logFields).Msgf("scheduler: %s scheduler failed", endpointName)
		} else {
			logger.Info().Fields(logFields).Msgf("scheduler: %s scheduler done", endpointName)
			middleware.ScheduleCallAddFuncAfterMiddleware(ctx)
		}

	}
	if spec == "" {
		logger.Info().Msgf("No Start scheduler: %s", endpointName)
		return
	}

	_, err := c.AddFunc(spec, fn)
	if err != nil {
		panic("AddFunc err: " + err.Error())
	}
	if strings.Contains(spec, "@") {
		go fn() // 避免被 一開始啟動就被 ch 卡住
	}
}
