package random

import (
	"math/rand"
	"time"
)

type RandomUtils interface {
	GenerateAlphanumericString(length int) string
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
