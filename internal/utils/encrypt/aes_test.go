package encrypt_test

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"testing"

	. "kedai/backend/be-kedai/internal/utils/encrypt"

	"github.com/stretchr/testify/assert"
)

func TestEncryptMessage(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	plaintext := "Hello, world!"
	ciphertext, err := EncryptMessage(plaintext, key)
	assert.NoError(t, err)

	// Decode the Base64-encoded ciphertext.
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	assert.NoError(t, err)

	// Extract the IV and ciphertext from the message.
	iv := decodedCiphertext[:aes.BlockSize]
	encryptedText := decodedCiphertext[aes.BlockSize:]

	// Make sure the plaintext was padded correctly.
	assert.Equal(t, len(encryptedText)%aes.BlockSize, 0)

	// Decrypt the ciphertext using the same key.
	decryptedText, err := DecryptMessage(ciphertext, key)
	assert.NoError(t, err)

	// Make sure the decrypted plaintext is the same as the original plaintext.
	assert.Equal(t, plaintext, decryptedText)

	// Make sure the IV is different each time.
	plaintext2 := "Goodbye, world!"
	ciphertext2, err := EncryptMessage(plaintext2, key)
	assert.NoError(t, err)
	decodedCiphertext2, err := base64.StdEncoding.DecodeString(ciphertext2)
	assert.NoError(t, err)
	iv2 := decodedCiphertext2[:aes.BlockSize]
	assert.False(t, bytes.Equal(iv, iv2))
}

func TestDecryptMessage(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	plaintext := "Hello, world!"
	ciphertext, err := EncryptMessage(plaintext, key)
	assert.NoError(t, err)

	// Decrypt the ciphertext using the same key.
	decryptedText, err := DecryptMessage(ciphertext, key)
	assert.NoError(t, err)

	// Make sure the decrypted plaintext is the same as the original plaintext.
	assert.Equal(t, plaintext, decryptedText)

	// Try to decrypt the ciphertext using a different key.
	key2 := make([]byte, 32)
	_, err = rand.Read(key2)
	assert.NoError(t, err)
	_, err = DecryptMessage(ciphertext, key2)
	assert.Error(t, err)

	// Try to decrypt an invalid ciphertext.
	invalidCiphertext := base64.StdEncoding.EncodeToString([]byte("invalid ciphertext"))
	_, err = DecryptMessage(invalidCiphertext, key)
	assert.Error(t, err)
}
