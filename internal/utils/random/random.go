package random

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type RandomUtils interface {
	GenerateAlphanumericString(length int) string
	GenerateNumericString(length int) string
	GenerateSecureUniqueToken() string
}

type randomUtilsImpl struct{}

type RandomUtilsConfig struct{}

func NewRandomUtils(cfg *RandomUtilsConfig) RandomUtils {
	return &randomUtilsImpl{}
}

func (u *randomUtilsImpl) GenerateAlphanumericString(length int) string {
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	alphaNum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = alphaNum[seededRand.Intn(len(alphaNum))]
	}

	return string(b)
}

func (u *randomUtilsImpl) GenerateNumericString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var result string
	for i := 0; i < length; i++ {
		digit := rand.Intn(10)
		result += strconv.Itoa(digit)
	}

	return result
}

func (u *randomUtilsImpl) GenerateSecureUniqueToken() string {
	b := make([]byte, 32)
	if _, err := cryptoRand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
