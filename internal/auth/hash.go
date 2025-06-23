package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	defaultCost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
