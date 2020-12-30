package util

// TODO: Tests schreiben
func ContainsString(haystack []string, needle string) bool {
	for _, currentNeedle := range haystack {
		if currentNeedle == needle {
			return true
		}
	}
	return false
}
