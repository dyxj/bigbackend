package account

import "golang.org/x/crypto/bcrypt"

const passwordHashCost = 14

func hashPassword(password string) (string, error) {
	hPass, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return "", err
	}
	return string(hPass), nil
}

// checkPasswordHash compares a plaintext password with a hashed password.
// returns nil if they match.
// returns ErrMismatchedHashAndPassword if they do not match.
// any other error should be handled as unexpected.
func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
