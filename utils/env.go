package utils

import (
	"os"
	"strings"
)

func MustGetenvStr(name, defaultValue string) string {
	env := os.Getenv(name)
	if StringEmptyOrBlank(env) {
		return defaultValue
	} else {
		return env
	}
}

func MustGetenvInt(name string, defaultValue int) int {
	env := os.Getenv(name)
	if StringEmptyOrBlank(env) {
		return defaultValue
	} else {
		if StringIsNumeric(env) {
			return Str2Int(env)
		} else {
			return defaultValue
		}
	}
}

func MustGetenvBool(name string, defaultValue bool) bool {
	env := os.Getenv(name)
	if !StringEmptyOrBlank(env) {
		if strings.EqualFold(env, "true") {
			return true
		}
		if strings.EqualFold(env, "false") {
			return false
		}
	}
	return defaultValue
}
