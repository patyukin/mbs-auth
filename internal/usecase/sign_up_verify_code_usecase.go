package usecase

import (
	"context"
)

func (u *UseCase) VerifySignUpCode(ctx context.Context, code string) error {
	panic("sdfsdf")
	//allSignUpCode, err := u.chr.GetSignUpCode(ctx, message.From.UserName)
	//if err != nil {
	//	log.Error().Msgf("Error getting sign up code: %v", err)
	//	msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid invite link.")
	//	if _, err = uc.bot.API.Send(msg); err != nil {
	//		log.Error().Msgf("Error sending message: %v", err)
	//		return
	//	}
	//
	//	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
	//	return
	//}
	//
	//result := strings.SplitN(allSignUpCode, ":", 2)
	//
	//signUpCode, err := uuid.Parse(result[0])
	//if err != nil {
	//	msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
	//	if _, err = uc.bot.API.Send(msg); err != nil {
	//		log.Error().Msgf("Error sending message: %v", err)
	//		return
	//	}
	//
	//	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
	//	return
	//}
	//
	//if signUpCode.String() != code.String() {
	//	log.Error().Msgf("Error parsing invite code: %v != %v", signUpCode.String(), code.String())
	//	msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid invite link.")
	//	if _, err = uc.bot.API.Send(msg); err != nil {
	//		log.Error().Msgf("Error sending message: %v", err)
	//		return
	//	}
	//
	//	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
	//	return
	//}
	//
	//err = uc.redis.DeleteSignUpCode(ctx, message.From.UserName)
	//if err != nil {
	//	log.Error().Msgf("Error deleting sign up code: %v", err)
	//	msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
	//	if _, err = uc.bot.API.Send(msg); err != nil {
	//		log.Error().Msgf("Error sending message: %v", err)
	//		return
	//	}
	//
	//	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
	//	return
	//}
	//
	//userUUID, err := uuid.Parse(result[1])
	//if err != nil {
	//	msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
	//	if _, err = uc.bot.API.Send(msg); err != nil {
	//		log.Error().Msgf("Error sending message: %v", err)
	//		return
	//	}
	//
	//	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
	//	return
	//}
	//
	//err = uc.registry.GetRepo().UpdateTelegramUserAfterSignUp(ctx, userUUID, message.Chat.ID, message.From.ID)
	//if err != nil {
	//	msg := tgbotapi.NewMessage(message.Chat.ID, "Error was encountered. Please try again.")
	//	if _, err = uc.bot.API.Send(msg); err != nil {
	//		log.Error().Msgf("Error sending message: %v", err)
	//		return
	//	}
	//
	//	log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
	//	return
	//}
	//
	//msg := tgbotapi.NewMessage(message.Chat.ID, "You have successfully signed up!")
	//if _, err = uc.bot.API.Send(msg); err != nil {
	//	log.Error().Msgf("Error sending message: %v", err)
	//	return
	//}
	//
	//log.Info().Msgf("Sent message to chat %d", message.Chat.ID)
}
