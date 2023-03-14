package hash

import (
	"crypto/hmac"
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
	key := []byte(config.HashKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(word))
	return hex.EncodeToString(h.Sum(nil))
}

func CompareSignature(apiSignature string, hashedSign string) bool {
	return hmac.Equal([]byte(apiSignature), []byte(hashedSign))
}
