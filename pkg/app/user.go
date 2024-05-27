package app

import (
	"fmt"
	"unicode"
)

func newUser(username string) (*User, error) {
	err := validateUsername(username)
	if err != nil {
		return nil, err
	}
	return &User{Username: username}, nil
}

type User struct {
	Username string `gorm:"column:username;type:varchar(64);not null;primaryKey"`
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
