package utils

func StringInArray(s string, arr []string) bool {
	for i := 0; i <= len(arr)-1; i++ {
		if s == arr[i] {
			return true
		}
	}
	return false
}
