package transfer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func parseTag(tag string) (isDnField, isTransient bool, attrName, dnAttr string, idx int) {
	tagParts := strings.Split(tag, ",")
	attrName = strings.TrimSpace(tagParts[0])

	for j := 1; j < len(tagParts); j++ {
		part := strings.TrimSpace(tagParts[j])
		if strings.HasPrefix(part, "dnAttr:") {
			dnAttr = strings.TrimPrefix(part, "dnAttr:")
		} else if strings.HasPrefix(part, "idx:") {
			if idxStr := strings.TrimPrefix(part, "idx:"); idxStr != "" {
				if parsedIdx, err := strconv.Atoi(idxStr); err == nil {
					idx = parsedIdx
				}
			}
		} else if part == "transient" {
			isTransient = true
		}
	}

	// 检查是否是 DN 字段
	if attrName == "dn" {
		isDnField = true
	}
	return
}

func ParseFromLdap[T any](entry *ldap.Entry) (*T, error) {
	var zero T

	// 检查是否为指针或结构体
	v := reflect.ValueOf(&zero).Elem()
	t := v.Type()

	// 确保是结构体
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ParseFromLdap: type %T is not a struct", zero)
	}

	item := &zero

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 跳过不可设置的字段（如未导出字段）
		if !fieldValue.CanSet() {
			continue
		}

		// 获取并解析 ldap tag
		ldapTag := field.Tag.Get("ldap")
		if ldapTag == "" || ldapTag == "-" {
			continue
		}

		isDnField, _, attrName, dnAttr, idx := parseTag(ldapTag)

		var attr string
		var values []string

		// 处理不同的字段类型
		if isDnField {
			// 处理 DN 字段
			if fieldValue.Kind() == reflect.String {
				// 如果字段类型是 string
				attr = entry.DN
			}
		} else if dnAttr != "" && idx >= 0 {
			// 从 DN 中提取特定属性
			dnParts := strings.Split(entry.DN, ",")
			for _, part := range dnParts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(strings.ToLower(part), strings.ToLower(dnAttr)+"=") {
					value := strings.SplitN(part, "=", 2)
					if len(value) == 2 {
						if idx == 0 || len(values) < idx {
							attr = strings.TrimSpace(value[1])
							break
						}
					}
				}
			}
		} else {
			// 从常规属性中获取值
			attr = entry.GetAttributeValue(attrName)
			values = entry.GetAttributeValues(attrName)
		}

		// 如果没有找到值，继续下一个字段
		if attr == "" && len(values) == 0 {
			continue
		}

		// 根据字段类型设置值
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(attr)
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.String {
				// 处理 []string 类型
				if len(values) == 0 && attr != "" {
					values = []string{attr}
				}
				slice := reflect.MakeSlice(fieldValue.Type(), len(values), len(values))
				for k, v := range values {
					slice.Index(k).SetString(v)
				}
				fieldValue.Set(slice)
			} else {
				return nil, fmt.Errorf("field %s: unsupported slice type %s", field.Name, fieldValue.Type().Elem().Kind())
			}
		default:
			return nil, fmt.Errorf("field %s: unsupported type %s", field.Name, fieldValue.Kind())
		}
	}

	return item, nil
}

func ParseToLdapAttributes[T any](item *T) (map[string][]string, error) {
	v := reflect.ValueOf(item).Elem()
	t := v.Type()

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ParseToLdapAttributes: type %T is not a struct", item)
	}

	attrs := make(map[string][]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 获取并解析 ldap tag
		ldapTag := field.Tag.Get("ldap")
		if ldapTag == "" || ldapTag == "-" {
			continue
		}

		isDnField, isTransient, attrName, _, _ := parseTag(ldapTag)

		// 跳过 transient 字段和 DN 字段（DN字段通常在创建条目时单独处理
		if isTransient || isDnField {
			continue
		}

		// 根据字段类型处理
		switch fieldValue.Kind() {
		case reflect.String:
			val := fieldValue.String()
			if val != "" {
				attrs[attrName] = []string{val}
			}
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.String {
				slice := fieldValue
				if slice.Len() > 0 {
					values := make([]string, slice.Len())
					for j := 0; j < slice.Len(); j++ {
						values[j] = slice.Index(j).String()
					}
					attrs[attrName] = values
				}
			} else {
				return nil, fmt.Errorf("field %s: unsupported slice type %s", field.Name, fieldValue.Type().Elem().Kind())
			}
		default:
			return nil, fmt.Errorf("field %s: unsupported type %s", field.Name, fieldValue.Kind())
		}
	}

	return attrs, nil
}
