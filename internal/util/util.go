package util

func ContainsString(haystack []string, needle string) bool {
	for _, b := range haystack {
		if b == needle {
			return true
		}
	}
	return false
}
