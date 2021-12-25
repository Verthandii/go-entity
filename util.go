package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func WriteFile(data []byte, output string, filename string) {
	_, err := os.Open(output)
	if err != nil {
		err = os.MkdirAll(output, 0644)
		if err != nil {
			log.Fatalln("os mkdir occurred error:", err)
		}
	}
	err = ioutil.WriteFile(output+filename, data, 0644)
	if err != nil {
		log.Fatalln("write file occurred error:", err)
	}
	log.Println(fmt.Sprintf("complete file %s%s", output, filename))
}
