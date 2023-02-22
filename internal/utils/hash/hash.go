package hash

import "golang.org/x/crypto/bcrypt"

func HashAndSalt(pwd string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	return string(hash), nil
}

func ComparePassword(hashedPw string, inputPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(inputPw))
	return err == nil
}
