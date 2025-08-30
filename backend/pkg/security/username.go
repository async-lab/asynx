package security

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func ValidateMemberUsernameLegality(username string) error {
	// 检查长度是否为10位
	if len(username) != 10 {
		return fmt.Errorf("用户名长度必须为10位，当前长度为%d位", len(username))
	}

	// 检查是否为纯数字
	pattern := `^\d+$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(username) {
		return fmt.Errorf("用户名必须为纯数字，当前值为%s", username)
	}

	// 检查前4位是否为合理的年份
	yearStr := username[:4]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return fmt.Errorf("用户名前4位必须为有效年份，当前值为%s: %w", yearStr, err)
	}

	// 获取当前年份
	currentYear := time.Now().Year()

	// 假设年份范围：从2000年到未来5年（可根据实际情况调整）
	minYear := 2000
	maxYear := currentYear + 5

	if year < minYear || year > maxYear {
		return fmt.Errorf("用户名前4位年份必须在%d到%d之间，当前年份为%d", minYear, maxYear, year)
	}

	return nil
}
