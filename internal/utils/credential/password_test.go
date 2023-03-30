package credential_test

import (
	"fmt"
	"kedai/backend/be-kedai/internal/utils/credential"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    bool
	}{
		{
			description: "should return false because there are no uppercase letter",
			input:       "password123",
			expected:    false,
		},
		{
			description: "should return false because there are no lowercase letter",
			input:       "PASSWORD123",
			expected:    false,
		},
		{
			description: "should return false because there are no number",
			input:       "Password",
			expected:    false,
		},
		{
			description: "should return false because there are emojis",
			input:       "Password123🫶",
			expected:    false,
		},
		{
			description: "should return true because password includes at least one uppercase, one lowercase, one number, and no emoji",
			input:       "Password123",
			expected:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			res := credential.VerifyPassword(tc.input)

			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestContainsUsername(t *testing.T) {
	type input struct {
		pw       string
		username string
	}
	cases := []struct {
		description string
		input       input
		expected    bool
	}{
		{
			description: "should return false because there are no username",
			input: input{
				pw:       "John12312",
				username: "notasd",
			},
			expected: false,
		},
		{
			description: "should return true because there are username",
			input: input{
				pw:       "John12312",
				username: "John",
			},
			expected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			res := credential.ContainsUsername(tc.input.pw, tc.input.username)

			assert.Equal(t, tc.expected, res)
		})
	}

}

func BenchmarkVerifyPassword(b *testing.B) {
	var table = []struct {
		name string
		pw   string
	}{
		{name: "valid", pw: "Password123"},
		{name: "invalid_length", pw: "pw1"},
		{name: "invalid_no_upper", pw: "password123"},
		{name: "invalid_no_lower", pw: "PASSWORD123"},
		{name: "invalid_no_numeric", pw: "Password"},
		{name: "invalid_contains_username", pw: "Passworduser123"},
		{name: "invalid_contains_emoji", pw: "😊Password123"},
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("password_"+v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				credential.VerifyPassword(v.pw)
			}
		})
	}
}

func BenchmarkContainsUsername(b *testing.B) {
	var table = []struct {
		name     string
		password string
		username string
	}{
		{name: "match", password: "myPassword", username: "pass"},
		{name: "no_match", password: "myPassword", username: "foo"},
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("contains_username_"+v.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				credential.ContainsUsername(v.password, v.username)
			}
		})
	}
}
