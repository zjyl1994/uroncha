package uroncha

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/zjyl1994/caasiu"
)

func init() {
	caasiu.RegisterRule("not_all", func(ruleName string, fieldName string, value interface{}) (bool, string) {
		replacement := strings.TrimPrefix(ruleName, "not_all:")
		if value != nil {
			if reflect.TypeOf(value).String() == "string" && len(strings.Replace(value.(string), replacement, "", -1)) == 0 {
				return false, fmt.Sprintf("The %s field must a string and not be empty", fieldName)
			}
		}
		return true, ""
	})

	caasiu.RegisterRule("min_numeric", func(ruleName string, fieldName string, value interface{}) (bool, string) {
		minNumber, _ := strconv.ParseInt(strings.TrimPrefix(ruleName, "min_numeric:"), 10, 64)
		val, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil || val < minNumber {
			return false, fmt.Sprintf("The %s field must be an int numeric and greate than %d or equal to", fieldName, minNumber)
		}
		return true, ""
	})

	caasiu.RegisterRule("max_numeric", func(ruleName string, fieldName string, value interface{}) (bool, string) {
		maxNumber, _ := strconv.ParseInt(strings.TrimPrefix(ruleName, "max_numeric:"), 10, 64)
		val, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil || val > maxNumber {
			return false, fmt.Sprintf("The %s field must be an int numeric and less than %d or equal to", fieldName, maxNumber)
		}
		return true, ""
	})

	caasiu.RegisterRule("utc_timestamp", func(ruleName string, fieldName string, value interface{}) (bool, string) {
		val := value.(string)
		msInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, fmt.Sprintf("The %s field must be numeric", fieldName)
		}
		// The value must be greater than utc-timestamp for 1970-01-01 08:00:12.133
		if msInt < 12133 {
			return false, fmt.Sprintf("The %s field must be right utc-timestamp(mill-second) and greater than 12133 (utc-timestamp for 1970-01-01 08:00:12.133)", fieldName)
		}
		return true, ""
	})

	caasiu.RegisterRule("unix_timestamp", func(ruleName string, fieldName string, value interface{}) (bool, string) {
		val := value.(string)
		msInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, fmt.Sprintf("The %s field must be numeric", fieldName)
		}
		// The value must be greater than unix-timestamp for 1970-01-01 08:00:12.133
		if msInt < 12 {
			return false, fmt.Sprintf("The %s field must be right unix-timestamp(second) and greater than 12 (unix-timestamp for 1970-01-01 08:00:12.133)", fieldName)
		}
		return true, ""
	})

	caasiu.RegisterRule("datetime", func(ruleName string, fieldName string, value interface{}) (bool, string) {
		const datetimeRegex = `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]) ([01]\d|2[0-3])(:[0-5]\d){2}$`
		if datetimeStr, ok := value.(string); ok {
			if match, _ := regexp.MatchString(datetimeRegex, datetimeStr); match {
				return true, ""
			}
		}
		return false, fmt.Sprintf("The %s field must be yyyy-MM-dd HH:mm:ss format", fieldName)
	})
}
