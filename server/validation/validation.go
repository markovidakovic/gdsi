package validation

import "regexp"

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsValidPhone(val string) bool {
	return regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,15}$`).MatchString(val)
}
