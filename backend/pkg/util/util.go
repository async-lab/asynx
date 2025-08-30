package util

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func GetAttributeKeys[T any]() ([]string, error) {
	var zero T
	v := reflect.ValueOf(&zero).Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("GetAttributeKeys: type %T is not a struct", zero)
	}

	seen := make(map[string]struct{}) // 用于去重
	attrs := make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		attrName := field.Tag.Get("ldap")
		if attrName != "" && attrName != "-" {
			if _, exists := seen[attrName]; !exists {
				seen[attrName] = struct{}{}
				attrs = append(attrs, attrName)
			}
		}
	}
	return attrs, nil
}

func BuildObjectClassCondition(objectClasses []string) string {
	conditions := make([]string, 0, len(objectClasses))
	for _, oc := range objectClasses {
		conditions = append(conditions, fmt.Sprintf("(objectClass=%s)", oc))
	}
	return fmt.Sprintf("(&%s)", strings.Join(conditions, ""))
}

// FindFirstMissingPositive 接收一个整数切片，返回从该切片最小值开始的第一个非连续的整数。
// 例如：
// [5, 6, 7, 9] -> 8
// [1, 2, 3, 4] -> 5
// [-3, -2, 0, 1] -> -1
// [] -> 1 (可以自定义，这里假设从1开始)
func FindFirstMissingPositive(nums []int) int {
	// 1. 处理空切片的特殊情况。
	// 如果列表为空，没有最小值作为起点，我们可以返回一个默认值，比如1。
	if len(nums) == 0 {
		return 1
	}

	// 2. 为了检查连续性，必须先对切片进行排序。
	sort.Ints(nums)

	// 3. 去除重复项，因为重复项会干扰连续性检查。
	// 例如 [5, 5, 7]，我们希望检查 5 和 7 之间的关系。
	uniqueNums := make([]int, 0, len(nums))
	if len(nums) > 0 {
		uniqueNums = append(uniqueNums, nums[0])
		for i := 1; i < len(nums); i++ {
			if nums[i] != nums[i-1] {
				uniqueNums = append(uniqueNums, nums[i])
			}
		}
	}

	// 4. 从第二个元素开始遍历去重后的切片。
	// 检查当前元素是否是前一个元素加一。
	for i := 1; i < len(uniqueNums); i++ {
		// 如果发现不连续，例如 [..., 5, 7, ...]
		// 那么 5+1=6 就是第一个空位。
		if uniqueNums[i] > uniqueNums[i-1]+1 {
			return uniqueNums[i-1] + 1
		}
	}

	// 5. 如果循环结束都没有找到空位，说明列表是连续的。
	// 那么第一个空位就是最大值（即排序后最后一个元素）加一。
	return uniqueNums[len(uniqueNums)-1] + 1
}
