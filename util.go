package main

import "strings"

// ToBigCamelCase user_id => UserId
func ToBigCamelCase(str string) string {
	result := ""
	if str == "" {
		return result
	}

	strs := strings.Split(str, "_")
	for _, s := range strs {
		r := []rune(s)
		if len(r) > 0 {
			if r[0] >= 'a' && r[0] <= 'z' {
				r[0] -= 32
			}
			result += string(r)
		}
	}

	return result
}

func MaxFunc(i, j int) int {
	if i > j {
		return i
	}
	return j
}
