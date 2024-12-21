package utils

import "strings"

func IsExist(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func IsContain(target string, slice []string) bool {
	for _, v := range slice {
		if strings.Contains(target, v) {
			return true
		}
	}
	return false
}

func SubString(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)

	if start < 0 {
		start = 0
	}

	if start >= rl {
		return ""
	}

	end := start + length
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
