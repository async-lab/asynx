package security

import (
	"fmt"
	"strings"
)

var LegalChars = "123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()-_=+[]{}|;:,.<>?/`~'"

func ValidatePasswordLegality(password string) error {
	if len(password) > 64 {
		return fmt.Errorf("password length exceeds 64 characters")
	}

	for _, ch := range password {
		if !strings.ContainsRune(LegalChars, ch) {
			return fmt.Errorf("password contains illegal characters: %q", ch)
		}
	}
	return nil
}

func ValidatePasswordStrength(password string) error {
	// var hasUpper, hasLower, hasDigit, hasSpecial bool
	// for _, ch := range password {
	// 	switch {
	// 	case 'A' <= ch && ch <= 'Z':
	// 		hasUpper = true
	// 	case 'a' <= ch && ch <= 'z':
	// 		hasLower = true
	// 	case '0' <= ch && ch <= '9':
	// 		hasDigit = true
	// 	case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:,.<>?/`~'", ch):
	// 		hasSpecial = true
	// 	}
	// }

	if len(password) < 12 {
		return fmt.Errorf("password is too short; it must be at least 12 characters long")
	}
	// if !hasUpper {
	// 	return fmt.Errorf("password must contain at least one uppercase letter")
	// }
	// if !hasLower {
	// 	return fmt.Errorf("password must contain at least one lowercase letter")
	// }
	// if !hasDigit {
	// 	return fmt.Errorf("password must contain at least one digit")
	// }
	// if !hasSpecial {
	// 	return fmt.Errorf("password must contain at least one special character")
	// }

	return nil
}
