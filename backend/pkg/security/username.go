package security

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func ValidateMemberUsernameLegality(username string) error {
	// 检查长度是否为10
	if len(username) != 10 {
		return fmt.Errorf("username must be 10 characters long, got %d", len(username))
	}

	// 检查是否全为数字
	pattern := `^\d+$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(username) {
		return fmt.Errorf("username must be all digits, got %s", username)
	}

	// 检查前4位是否为有效年份
	yearStr := username[:4]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return fmt.Errorf("the first 4 characters must be a valid year, got %s: %w", yearStr, err)
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	// 假设年份范围：2000年至当前年份后5年
	minYear := 2000
	maxYear := currentYear + 5

	if year < minYear || year > maxYear {
		return fmt.Errorf("the year in the first 4 characters must be between %d and %d, got %d", minYear, maxYear, year)
	}

	return nil
}
