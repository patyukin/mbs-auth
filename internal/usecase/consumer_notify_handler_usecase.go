package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/mbs-pkg/pkg/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"strings"
)

func (u *UseCase) NotifyConsumeHandler(ctx context.Context, msg amqp.Delivery) error {
	var message model.AuthSignUpConfirmCode
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	allSignUpCode, err := u.chr.GetSignUpCode(ctx, message.UserTelegramLogin)
	log.Debug().Msgf("allSignUpCode: %s", allSignUpCode)
	if err != nil {
		err = sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Invalid invite link. Please try again.")
		if err != nil {
			return fmt.Errorf("failed to send message to queue: %w", err)
		}

		return fmt.Errorf("failed to get sign up code: %w", err)
	}

	result := strings.SplitN(allSignUpCode, ":", 2)

	signUpCode, err := uuid.Parse(result[0])
	if err != nil {
		err = sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Error was encountered. Please try again.")
		if err != nil {
			return fmt.Errorf("failed to send message to queue: %w", err)
		}

		return fmt.Errorf("failed to parse sign up code: %w", err)
	}

	if signUpCode.String() != message.Code {
		err = sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Error was encountered. Please try again.")
		if err != nil {
			return fmt.Errorf("failed to send message to queue: %w", err)
		}

		return fmt.Errorf("invalid sign up code: %w", err)
	}

	err = u.chr.DeleteSignUpCode(ctx, message.UserTelegramLogin)
	if err != nil {
		err = sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Error was encountered. Please try again.")
		if err != nil {
			return fmt.Errorf("failed to send message to queue: %w", err)
		}

		return fmt.Errorf("failed to delete sign up code: %w", err)
	}

	userUUID, err := uuid.Parse(result[1])
	if err != nil {
		return sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Error was encountered. Please try again.")
	}

	err = u.registry.GetRepo().UpdateTelegramUserAfterSignUp(ctx, userUUID, message.ChatID, message.UserTelegramID)
	if err != nil {
		return sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Error was encountered. Please try again.")
	}

	return sendToQueue(ctx, u.prdcr.PublishAuthSignUpResultMessage, message.ChatID, "Success!")
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
