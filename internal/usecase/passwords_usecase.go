package usecase

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const CostSalt = 12

func (u *UseCase) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), CostSalt)
	if err != nil {
		return "", fmt.Errorf("failed bcrypt.GenerateFromPassword: %w", err)
	}

	return string(hashedPassword), nil
}

func (u *UseCase) ComparePasswords(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return fmt.Errorf("failed bcrypt.CompareHashAndPassword: %w", err)
	}

	return nil
}
