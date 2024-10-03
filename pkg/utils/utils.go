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
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
func ParseTime(timeStr string) (*time.Time, error) {
	layout := "2006-01-02T15:04-07:00"

	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	return &parsedTime, nil
}

func ParseUUID(uuidStr string) *uuid.UUID {
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil
	}

	return &uuid
}
