package usecase

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (u *UseCase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed bcrypt.GenerateFromPassword: %w", err)
	}

	return string(hashedPassword), nil
}

func (u *UseCase) ComparePasswords(hashedPassword []byte, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plainPassword))
	return err
}

func (u *UseCase) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Error().Msgf("failed bcrypt.CompareHashAndPassword: %v", err)
		return false
	}

	return true
}
