package password_test

import (
	"kedai/backend/be-kedai/internal/utils/password"
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
			res := password.VerifyPassword(tc.input)

			assert.Equal(t, tc.expected, res)
		})
	}
}
