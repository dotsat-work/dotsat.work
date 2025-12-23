package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	MinPasswordLength   = 12
	MaxPasswordLength   = 72 // bcrypt limit
	MaxConsecutiveChars = 6
)

func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}

	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password must not exceed %d characters", MaxPasswordLength)
	}

	lower := strings.ToLower(password)

	// High-risk passwords only
	commonPasswords := []string{
		"password", "passw0rd",
		"123456", "12345678", "123456789", "1234567890",
		"qwerty", "qwertyuiop", "asdfghjkl",
		"letmein", "admin", "welcome", "monkey",
	}

	for _, weak := range commonPasswords {
		// Exact match
		if lower == weak {
			return errors.New("password is too common")
		}

		// weak + numbers/special OR numbers/special + weak
		pattern := `[\d!@#$%^&*]+`
		if regexp.MustCompile(`^`+weak+pattern+`$`).MatchString(lower) ||
			regexp.MustCompile(`^`+pattern+weak+`$`).MatchString(lower) {
			return errors.New("password is too common")
		}
	}

	// Reject excessive repetition (6+ same characters)
	if hasExcessiveRepetition(password) {
		return errors.New("password contains excessive repetition")
	}

	return nil
}

// hasExcessiveRepetition checks if the password has MaxConsecutiveChars or more consecutive identical characters
// Note: No length check needed - password length is already validated in ValidatePassword
func hasExcessiveRepetition(password string) bool {
	count := 1
	for i := 1; i < len(password); i++ {
		if password[i] == password[i-1] {
			count++
			if count >= MaxConsecutiveChars {
				return true
			}
		} else {
			count = 1
		}
	}
	return false
}
