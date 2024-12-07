package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/patyukin/mbs-auth/internal/db"
	"github.com/patyukin/mbs-pkg/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (u *UseCase) RegistrationSolutionProcess(ctx context.Context, record *kgo.Record) error {
	var message model.AuthSignUpConfirmCode

	log.Debug().Msgf("Received record: %v", string(record.Value))

	if err := json.Unmarshal(record.Value, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	err := u.registry.ReadCommitted(
		ctx, func(ctx context.Context, repo *db.Repository) error {
			allSignUpCode, err := u.chr.GetSignUpCode(ctx, message.UserTelegramLogin)
			log.Debug().Msgf("allSignUpCode: %s", allSignUpCode)
			if err != nil {
				return fmt.Errorf("failed to get sign up code: %w", err)
			}

			result := strings.SplitN(allSignUpCode, ":", 2)

			signUpCode, err := uuid.Parse(result[0])
			if err != nil {
				return fmt.Errorf("failed to parse sign up code: %w", err)
			}

			if signUpCode.String() != message.Code {
				return fmt.Errorf("invalid sign up code: %w", err)
			}

			err = u.chr.DeleteSignUpCode(ctx, message.UserTelegramLogin)
			if err != nil {
				return fmt.Errorf("failed to delete sign up code: %w", err)
			}

			userUUID, err := uuid.Parse(result[1])
			if err != nil {
				return fmt.Errorf("failed to parse user uuid: %w", err)
			}

			err = repo.UpdateTelegramUserAfterSignUp(ctx, userUUID, message.ChatID, message.UserTelegramID)
			if err != nil {
				return fmt.Errorf("failed to update telegram user: %w", err)
			}

			// add user role
			_, err = repo.AddUserToRole(ctx, userUUID.String(), "user")
			if err != nil {
				return fmt.Errorf("failed to add user to role: %w", err)
			}

			return nil
		},
	)
	if err != nil {
		sendErr := sendToQueue(ctx, u.prdcr.EnqueueTelegramMessage, message.ChatID, "Error was encountered. Please try again.")
		if sendErr != nil {
			return fmt.Errorf("failed to execute transaction, err: %w, failed to send message to queue: %w", err, sendErr)
		}

		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return sendToQueue(ctx, u.prdcr.EnqueueTelegramMessage, message.ChatID, "Success.")
}

func sendToQueue(ctx context.Context, handler func(ctx context.Context, body []byte, headers amqp.Table) error, chatID int64, msg string) error {
	result := model.AuthSignUpResultMessage{
		ChatID:  chatID,
		Message: msg,
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	if err = handler(ctx, resultBytes, amqp.Table{}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
