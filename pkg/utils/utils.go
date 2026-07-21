package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPasswordStr := string(hashedPassword)

	return &hashedPasswordStr, nil
}

func ComparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func ParseTime(timeStr string) (*time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	return &parsedTime, nil
}

func ParseUUID(uuidStr string) *uuid.UUID {
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil
	}

	return &parsedUUID
}
