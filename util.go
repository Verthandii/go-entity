package main

import (
	"strings"
)

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
			for _, word := range Cfg.BigCamelCaseWords {
				if word == s {
					for i := range r {
						r[i] -= 32
					}
				}
			}
			if r[0] >= 'a' && r[0] <= 'z' {
				r[0] -= 32
			}
			result += string(r)
		}
	}

	return result
}

// ToSmallCamelCase user_id => userId
func ToSmallCamelCase(str string) string {
	result := ""
	if str == "" {
		return result
	}

	strs := strings.Split(str, "_")
	for i, s := range strs {
		r := []rune(s)
		if len(r) > 0 {
			if i == 0 {
				result += string(r)
				continue
			}

			for _, word := range Cfg.BigCamelCaseWords {
				if word == s {
					for j := range r {
						r[j] -= 32
					}
				}
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
