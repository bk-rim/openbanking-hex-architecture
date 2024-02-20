package helper

import "regexp"

func IsValidIban(iban string) bool {
	ibanRegex := regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]{1,30}$`)
	return ibanRegex.MatchString(iban)
}
