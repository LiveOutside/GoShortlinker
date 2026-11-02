// Package password содержит единые функции хеширования и проверки паролей.
package password

import "golang.org/x/crypto/bcrypt"

const cost = 12

func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
	return string(b), err
}

func Check(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
