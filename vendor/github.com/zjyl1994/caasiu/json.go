package caasiu

import (
	"fmt"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
)

type JSON struct {
	rawJSON []byte
}

func NewJSON(json []byte) *JSON {
	var result JSON
	result.rawJSON = json
	return &result
}

func (j *JSON) Valid(rule map[string][]string) []string {
	sjson, err := simplejson.NewJson(j.rawJSON)
	if err != nil {
		return []string{err.Error()}
	}
	var errMsg []string
	for fieldName, rulesOnField := range rule {
		if len(rulesOnField) > 0 {
			fieldPaths := strings.Split(fieldName, ".")
			currentJSONLevel := sjson.GetPath(fieldPaths...)
			if currentJSONLevel.Interface() == nil {
				if stringInArray("required", rulesOnField) {
					errMsg = append(errMsg, fmt.Sprintf(`field "%s" is required`, fieldName))
				}
				continue
			}
			for _, oneRule := range rulesOnField {
				ruleCommand := strings.Split(oneRule, ":")
				if ruleFunc, ok := registerRules[ruleCommand[0]]; ok {
					valid, errMessage := ruleFunc(oneRule, fieldName, currentJSONLevel.Interface())
					if !valid {
						errMsg = append(errMsg, errMessage)
					}
				}
			}
		}
	}
	return errMsg
}

func (j *JSON) Data() *simplejson.Json {
	sjson, _ := simplejson.NewJson(j.rawJSON)
	return sjson
}

func (j *JSON) Raw() []byte {
	return j.rawJSON
}
