package utils

import "os"

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
