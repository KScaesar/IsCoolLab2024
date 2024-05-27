package app

import (
	"fmt"
	"unicode"
)

func NewUser(username string) (*User, error) {
	return &User{Username: username}, nil
}

type User struct {
	Username     string
	FilesystemId string
}

func validateUsername(username string) error {
	if len(username) > 64 {
		return fmt.Errorf("Error: The %v %w", username, ErrInvalidParams)
	}
	for _, char := range username {
		if !(unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_' || char == '-') {
			return fmt.Errorf("Error: The %v %w", username, ErrInvalidParams)
		}
	}
	return nil
}
