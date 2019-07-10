package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func StringInArray(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func StringEmptyOrBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func StringPadLeft(sourceStr string, totalLen int, char string) string {
	if len(sourceStr) >= totalLen {
		return sourceStr
	} else {
		targetStrArr := []string{}
		for i := 0; i < totalLen-len(sourceStr); i++ {
			targetStrArr = append(targetStrArr, char)
		}
		targetStrArr = append(targetStrArr, sourceStr)
		return strings.Join(targetStrArr, "")
	}
}

func ToString(v interface{}) string {
	str, ok := v.(string)
	if !ok {
		str = fmt.Sprintf("%#v", v)
	}
	return str
}

func StringIsNumeric(str string) bool {
	return regexp.MustCompile("^[0-9]+$").MatchString(str)
}
