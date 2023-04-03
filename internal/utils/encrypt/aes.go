package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Encrypts the given plaintext message using the given key.
func EncryptMessage(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random initialization vector (IV).
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Pad the plaintext message to a multiple of the block size.
	paddedPlaintext := pkcs7Pad([]byte(plaintext), aes.BlockSize)

	// Create a new AES cipher block mode with the IV.
	ciphertext := make([]byte, len(paddedPlaintext))
	stream := cipher.NewCTR(block, iv)

	// Encrypt the padded plaintext message.
	stream.XORKeyStream(ciphertext, paddedPlaintext)

	// Concatenate the IV and ciphertext into a single message.
	message := make([]byte, aes.BlockSize+len(ciphertext))
	copy(message, iv)
	copy(message[aes.BlockSize:], ciphertext)

	// Encode the message as a Base64 string.
	return base64.StdEncoding.EncodeToString(message), nil
}

// Decrypts the given ciphertext message using the given key.
func DecryptMessage(ciphertext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Decode the Base64-encoded message.
	message, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Extract the IV and ciphertext from the message.
	fmt.Println(aes.BlockSize, len(message))
	if len(message)%aes.BlockSize != 0 {
		return "", errors.New("invalid encoded string")
	}
	iv := message[:aes.BlockSize]
	ciphertext = string(message[aes.BlockSize:])

	// Create a new AES cipher block mode with the IV.
	paddedPlaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)

	// Decrypt the ciphertext message.
	stream.XORKeyStream(paddedPlaintext, []byte(ciphertext))

	// Unpad the decrypted message.
	plaintext, err := pkcs7Unpad(paddedPlaintext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Pads the given input to a multiple of the given block size using the PKCS#7 scheme.
func pkcs7Pad(input []byte, blockSize int) []byte {
	padding := blockSize - len(input)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(input, padText...)
}

// Unpads the given input using the PKCS#7 scheme.
func pkcs7Unpad(input []byte) ([]byte, error) {
	length := len(input)
	padding := int(input[length-1])

	if padding > length {
		return nil, fmt.Errorf("invalid padding")
	}

	return input[:length-padding], nil
}
