package hash_test

import (
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
