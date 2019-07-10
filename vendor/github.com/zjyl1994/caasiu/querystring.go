package caasiu

import (
	"fmt"
	"net/url"
	"strings"
)

type QueryString struct {
	raw  url.Values
	data map[string]string
}

func NewQueryString(value url.Values) *QueryString {
	var result QueryString
	result.raw = value
	// preprocess
	result.data = make(map[string]string)
	for k, v := range value {
		if strings.TrimSpace(v[0]) != "" {
			result.data[k] = v[0]
		}
	}
	return &result
}

func (q *QueryString) Valid(rule map[string][]string) []string {
	var errMsg []string
	for fieldName, rulesOnField := range rule {
		if len(rulesOnField) > 0 {
			v, ok := q.data[fieldName]
			if ok {
				for _, oneRule := range rulesOnField {
					ruleCommand := strings.Split(oneRule, ":")
					if ruleFunc, ok := registerRules[ruleCommand[0]]; ok {
						valid, errMessage := ruleFunc(oneRule, fieldName, v)
						if !valid {
							errMsg = append(errMsg, errMessage)
						}
					}
				}
			} else {
				if stringInArray("required", rulesOnField) {
					errMsg = append(errMsg, fmt.Sprintf(`field "%s" is required`, fieldName))
				}
			}
		}
	}
	return errMsg
}

func (q *QueryString) Data() map[string]string {
	return q.data
}

func (q *QueryString) Raw() url.Values {
	return q.raw
}
