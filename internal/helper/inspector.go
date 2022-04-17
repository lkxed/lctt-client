package helper

func StringSliceContains(slice []string, target any) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}
