package utils

// CheckNA is sugar for build variables output.
func CheckNA(str string) string {
	if str == "" {
		return "N/A"
	}
	return str
}
