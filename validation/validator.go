package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername  = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidUFullname = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain characters from %d to %d", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters, digit, or underscore")
	}

	return nil
}

func ValidateFullname(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUFullname(value) {
		return fmt.Errorf("must contain only letters or spacce")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not valid email format")
	}

	return nil
}
