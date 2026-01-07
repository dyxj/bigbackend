package validx

import (
	"net/mail"
)

func IsEmail(email string) bool {
	address, err := mail.ParseAddress(email)
	if err == nil && address.Address == email {
		return true
	}
	return false
}
