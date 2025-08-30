package security

import (
	"fmt"
	"regexp"
)

// ValidateMailFormat 验证邮箱格式是否正确
func ValidateEmailFormat(email string) error {
	if email == "" {
		return fmt.Errorf("mail cannot be empty")
	}

	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if re.MatchString(email) {
		return nil
	} else {
		return fmt.Errorf("mail format is incorrect")
	}
}
