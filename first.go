package golibrary

import "strings"

func FirstString(str ...string) string {
	for i := range str {
		if strings.TrimSpace(str[i]) != "" {
			return str[i]
		}
	}
	return ""
}

func FirstInt64(items ...int64) int64 {
	for i := range items {
		if items[i] > 0 {
			return items[i]
		}
	}
	return 0
}

func FirstError(rst ...error) error {
	for _, v := range rst {
		if v != nil {
			return v
		}
	}
	return nil
}
