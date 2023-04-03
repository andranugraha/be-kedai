package hash_test

import (
	"fmt"
	"testing"

	"kedai/backend/be-kedai/internal/utils/hash"

	"github.com/stretchr/testify/assert"
)

func TestComparePassword(t *testing.T) {
	type input struct {
		actualPassword string
		inputPassword  string
	}

	cases := []struct {
		desciption string
		input
		expected bool
	}{
		{
			desciption: "should return true if inputted password is valid",
			input: input{
				actualPassword: "password",
				inputPassword:  "password",
			},
			expected: true,
		},
		{
			desciption: "should return false if inputted password is invalid",
			input: input{
				actualPassword: "password",
				inputPassword:  "another password",
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		hashedPassword, _ := hash.HashAndSalt(tc.input.actualPassword)

		res := hash.ComparePassword(hashedPassword, tc.input.inputPassword)

		assert.Equal(t, tc.expected, res)
	}
}

func BenchmarkHashAndSalt(b *testing.B) {
	var table = []struct {
		name string
		pwd  string
	}{
		{name: "short_pwd", pwd: "1234"},
		{name: "medium_pwd", pwd: "testpassword1234"},
		{name: "long_pwd", pwd: "This is a very long password that is more than 30 characters"},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("input_type_%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = hash.HashAndSalt(v.pwd)
			}
		})
	}
}

func BenchmarkComparePassword(b *testing.B) {
	hashed, _ := hash.HashAndSalt("password123") // Hash a sample password to use in the benchmark
	var table = []struct {
		name     string
		password string
	}{
		{name: "correct_password", password: "password123"},
		{name: "incorrect_password", password: "password456"},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("input_type_%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hash.ComparePassword(hashed, v.password)
			}
		})
	}
}

func BenchmarkHashSHA256(b *testing.B) {
	var table = []struct {
		name string
		word string
	}{
		{name: "short_word", word: "hello"},
		{name: "medium_word", word: "thequickbrownfoxjumpsoverthelazydog"},
		{name: "long_word", word: "This is a very long word that is more than 30 characters"},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("input_type_%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hash.HashSHA256(v.word)
			}
		})
	}
}

func BenchmarkCompareSignature(b *testing.B) {
	var table = []struct {
		name         string
		apiSignature string
		hashedSign   string
	}{
		{name: "match", apiSignature: "ThisIsASecretSignature", hashedSign: hash.HashSHA256("ThisIsASecretSignature")},
		{name: "mismatch", apiSignature: "ThisIsASecretSignature", hashedSign: hash.HashSHA256("ThisIsNotTheSameSignature")},
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("input_type_%s", v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				hash.CompareSignature(v.apiSignature, v.hashedSign)
			}
		})
	}
}
