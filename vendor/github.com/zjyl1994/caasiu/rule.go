package caasiu

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	regexInteger = "^\\d+$"
)

var (
	registerRules map[string]func(string, string, interface{}) (bool, string)

	ErrAlreadyRegistered = errors.New("rule already registered")
)

func init() {
	registerRules = map[string]func(string, string, interface{}) (bool, string){
		"string":  ruleString,
		"integer": ruleInteger,
		"regex":   ruleRegexp,
		"in":      ruleIn,
		"bool":    ruleBool,
		"array":   ruleArray,

		"ascii":     ruleASCII,
		"alpha":     ruleAlpha,
		"numeric":   ruleNumeric,
		"alphanum":  ruleAlphaNum,
		"hexstring": ruleHexString,
		"printable": rulePrintableASCII,
		"base64":    ruleBase64,
	}
}

func RegisterRule(ruleName string, ruleFunc func(string, string, interface{}) (bool, string)) error {
	if _, ok := registerRules[ruleName]; ok {
		return ErrAlreadyRegistered
	} else {
		registerRules[ruleName] = ruleFunc
		return nil
	}
}

func ruleString(ruleName string, fieldName string, value interface{}) (bool, string) {
	_, ok := value.(string)
	if ok {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" must string`, fieldName)
	}
}

func basicRegexpRule(ruleName string, fieldName string, value interface{}, regex string, errMsg string) (bool, string) {
	valueString, ok := value.(string)
	if ok {
		if regexp.MustCompile(regex).MatchString(valueString) {
			return true, ""
		} else {
			return false, errMsg
		}
	} else {
		return false, fmt.Sprintf(`field "%s" must string`, fieldName)
	}
}

func ruleInteger(ruleName string, fieldName string, value interface{}) (bool, string) {
	errMsg := fmt.Sprintf(`field "%s" must integer`, fieldName)
	var valueString string
	switch t := value.(type) {
	case string:
		valueString = t
	case json.Number:
		valueString = t.String()
	default:
		return false, errMsg
	}
	if regexp.MustCompile(regexInteger).MatchString(valueString) {
		return true, ""
	} else {
		return false, errMsg
	}
}

func ruleRegexp(ruleName string, fieldName string, value interface{}) (bool, string) {
	ruleParamStartAt := strings.Index(ruleName, ":")
	if ruleParamStartAt == -1 {
		return false, "no param available"
	}
	param := ruleName[ruleParamStartAt+1:]
	var valueString string
	switch t := value.(type) {
	case string:
		valueString = t
	case json.Number:
		valueString = t.String()
	default:
		return false, fmt.Sprintf(`field "%s" can't cast to string`, fieldName)
	}
	regex, err := regexp.Compile(param)
	if err != nil {
		return false, "illegal regexp"
	}
	if regex.MatchString(valueString) {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" not match "%s"`, fieldName, param)
	}
}

func ruleIn(ruleName string, fieldName string, value interface{}) (bool, string) {
	ruleParamStartAt := strings.Index(ruleName, ":")
	if ruleParamStartAt == -1 {
		return false, "no param available"
	}
	param := ruleName[ruleParamStartAt+1:]
	params := strings.Split(param, ",")
	var valueString string
	switch t := value.(type) {
	case string:
		valueString = t
	case json.Number:
		valueString = t.String()
	default:
		return false, fmt.Sprintf(`field "%s" can't cast to string`, fieldName)
	}
	if stringInArray(valueString, params) {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" not in [%s]`, fieldName, param)
	}
}

func ruleBool(ruleName string, fieldName string, value interface{}) (bool, string) {
	if _, ok := value.(bool); ok {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" must be bool`, fieldName)
	}
}

func ruleArray(ruleName string, fieldName string, value interface{}) (bool, string) {
	if _, ok := value.([]interface{}); ok {
		return true, ""
	} else {
		return false, fmt.Sprintf(`field "%s" must be array`, fieldName)
	}
}

func ruleASCII(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^[\x00-\x7F]+$", fmt.Sprintf(`field "%s" must ASCII string`, fieldName))
}

func ruleAlpha(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^[a-zA-Z]+$", fmt.Sprintf(`field "%s" must alpha string`, fieldName))
}

func ruleAlphaNum(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^[0-9a-zA-Z]+$", fmt.Sprintf(`field "%s" must alphanum string`, fieldName))
}

func ruleNumeric(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^[0-9]+$", fmt.Sprintf(`field "%s" must numeric string`, fieldName))
}

func ruleHexString(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^[0-9a-fA-F]+$", fmt.Sprintf(`field "%s" must hex string`, fieldName))
}

func rulePrintableASCII(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^[\x20-\x7E]+$", fmt.Sprintf(`field "%s" must printable ASCII string`, fieldName))
}

func ruleBase64(ruleName string, fieldName string, value interface{}) (bool, string) {
	return basicRegexpRule(ruleName, fieldName, value, "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$", fmt.Sprintf(`field "%s" must base64 string`, fieldName))
}
