package cronjob

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type CronJob struct {
	c  *cron.Cron
	uc UseCase
}

type UseCase interface {
	CleanExpiredTokens(ctx context.Context) error
	RemoveNotRegisteredUsers(ctx context.Context) error
}

func New(uc UseCase) *CronJob {
	return &CronJob{
		c:  cron.New(),
		uc: uc,
	}
}

func (cj *CronJob) Stop() {
	cj.c.Stop()
}

func (cj *CronJob) Run(ctx context.Context) error {
	_, err := cj.c.AddFunc("*/10 * * * *", func() {
		log.Info().Msg("run cj.uc.CleanExpiredTokens")

		if localErr := cj.uc.CleanExpiredTokens(ctx); localErr != nil {
			log.Error().Msgf("failed cj.uc.CleanExpiredTokens, err: %v", localErr)
		}
	})
	if err != nil {
		return fmt.Errorf("failed adding cron job cj.uc.CleanExpiredTokens: %w", err)
	}

	_, err = cj.c.AddFunc("0 * * * *", func() {
		log.Info().Msg("run cj.uc.RemoveNotRegisteredUsers")

		if localErr := cj.uc.RemoveNotRegisteredUsers(ctx); localErr != nil {
			log.Error().Msgf("failed cj.uc.RemoveNotRegisteredUsers, err: %v", localErr)
		}
	})
	if err != nil {
		return fmt.Errorf("failed adding cron job cj.uc.RemoveNotRegisteredUsers: %w", err)
	}

	cj.c.Start()
	return nil
}

// https://crontab.guru/
