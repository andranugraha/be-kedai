package credential_test

import (
	"kedai/backend/be-kedai/internal/utils/credential"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyUsername(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    bool
	}{
		{
			description: "should return false when username contains emoji",
			input:       "new_username🫶",
			expected:    false,
		},
		{
			description: "should return false when username doesn't contain any letter",
			input:       "127.0.0.1",
			expected:    false,
		},
		{
			description: "should return false when username contains any symbol other than '_' or '.'",
			input:       "n3w_u$er.name",
			expected:    false,
		},
		{
			description: "should return true when username contains at least one letter",
			input:       "new_u5ername",
			expected:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			verdict := credential.VerifyUsername(tc.input)

			assert.Equal(t, tc.expected, verdict)
		})
	}
}

func BenchmarkVerifyUsername(b *testing.B) {
	testCases := []struct {
		name     string
		username string
	}{
		{"valid_username", "john_doe.123"},
		{"invalid_username_contains_special_char", "john_doe!123"},
		{"invalid_username_contains_space", "john doe"},
		{"invalid_username_contains_emoji", "john_doe😀"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				credential.VerifyUsername(tc.username)
			}
		})
	}
}
