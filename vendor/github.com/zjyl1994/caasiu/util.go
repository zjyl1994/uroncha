package caasiu

func stringInArray(str string, arr []string) bool {
	for _, b := range arr {
		if b == str {
			return true
		}
	}
	return false
}
