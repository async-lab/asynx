package security

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var regexpDigitsOnly = regexp.MustCompile(`^\d+$`)

func ValidateMemberUsernameLegality(username string) error {
	// 1. 长度必须为10
	if len(username) != 10 {
		return fmt.Errorf("username must be 10 characters long, got %d", len(username))
	}

	// 2. 必须全为数字
	if !regexpDigitsOnly.MatchString(username) {
		return fmt.Errorf("username must be all digits, got %s", username)
	}

	currentYear := time.Now().Year()
	minYear, maxYear := 2000, currentYear+5

	var year int

	// 研究生学号：3 + 2位年份 + 7位其他数字
	if username[0] == '3' {
		shortYearStr := username[1:3] // 取第2、3位
		shortYear, err := strconv.Atoi(shortYearStr)
		if err != nil {
			return fmt.Errorf("the 2nd-3rd characters must be a valid 2-digit year, got %s: %w", shortYearStr, err)
		}
		year = 2000 + shortYear
	} else {
		yearStr := username[:4] // 普通学号：前4位是年份
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			return fmt.Errorf("the first 4 characters must be a valid year, got %s: %w", yearStr, err)
		}
		year = y
	}

	if year < minYear || year > maxYear {
		return fmt.Errorf("the year must be between %d and %d, got %d", minYear, maxYear, year)
	}

	return nil
}
