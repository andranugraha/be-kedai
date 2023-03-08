package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"kedai/backend/be-kedai/config"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) (string, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	return string(hash), nil
}

func ComparePassword(hashedPw string, inputPw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(inputPw))
	return err == nil
}

func HashSHA256(word string) string {
	h := sha256.New()
	h.Write([]byte(word + config.HashKey))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func CompareSignature(apiSignature string, comparedSignature string) bool {
	return apiSignature == comparedSignature
}
