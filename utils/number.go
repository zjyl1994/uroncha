package utils

import (
	"fmt"
	"strconv"
)

func Round(value float64, digital int) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(digital)+"f", value), 64)
	return value
}

func Str2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return i
	}
}

func Str2Int64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	} else {
		return i
	}
}

func Int2Str(i int) string {
	return strconv.Itoa(i)
}

func Int642Str(i int64) string {
	return strconv.FormatInt(i, 10)
}
