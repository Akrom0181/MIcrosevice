package etc

import "golang.org/x/crypto/bcrypt"

func GeneratePasswordHash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), 10)
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
