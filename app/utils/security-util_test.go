package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncryptDecryptAES128ECB tests both encryption and decryption for AES-128 in ECB mode.
func TestEncryptDecryptAES128ECB(t *testing.T) {
	// Sample plaintext and AES key
	plaintext := []byte("This is a secret message!")
	key := "1234567890123456" // 16-byte AES key (AES-128)

	// Encrypt the plaintext
	encrypted, err := EncryptAES128ECB(plaintext, key)
	if err != nil {
		log.Fatal("Encryption error:", err)
	}

	// Ensure the encrypted text is not equal to the original plaintext
	assert.NotEqual(t, string(plaintext), encrypted, "Encrypted text should not be the same as plaintext")

	// Decrypt the encrypted ciphertext
	decrypted, err := DecryptAES128ECB(encrypted, key)
	if err != nil {
		log.Fatal("Decryption error:", err)
	}

	// Ensure the decrypted text matches the original plaintext
	assert.Equal(t, string(plaintext), string(decrypted), "Decrypted text should match the original plaintext")
}

// TestEncryptAES128ECB tests encryption function independently
func TestEncryptAES128ECB(t *testing.T) {
	// Sample plaintext and AES key
	plaintext := []byte("Test message for encryption")
	key := "1234567890123456" // 16-byte AES key (AES-128)

	// Encrypt the plaintext
	encrypted, err := EncryptAES128ECB(plaintext, key)

	assert.NoError(t, err)

	// Ensure the result is a base64-encoded string and not empty
	assert.NotEmpty(t, encrypted, "Encrypted text should not be empty")

	// Further, we can check if the encrypted text length is greater than the plaintext length
	assert.Greater(t, len(encrypted), len(plaintext), "Encrypted text should be larger than plaintext")

	// Encrypt the plaintext
	_, err = EncryptAES128ECB(plaintext, "shortkey")

	assert.Error(t, err)
}

// Test for successful decryption with valid key and ciphertext.
func TestDecryptAES128ECB(t *testing.T) {
	enc, _ := EncryptAES128ECB([]byte("This is a secret"), "1234567890123456")
	log.Println("Encrypted: ", enc)

	tests := []struct {
		name          string
		ciphertext    string
		key           string
		expectedPlain string
		expectError   bool
	}{
		{
			name:          "Valid decryption",
			ciphertext:    "t8vP+hOjSj3a/nK0gzk+IgUBh6DN5amHLLqwkatz5VM=", // Example base64-encoded ciphertext
			key:           "1234567890123456",                             // 16-byte AES key
			expectedPlain: "This is a secret",                             // Expected plaintext after decryption
			expectError:   false,
		},
		{
			name:          "Invalid key",
			ciphertext:    "t8vP+hOjSj3a/nK0gzk+IgUBh6DN5amHLLqwkatz5VM=", // Example ciphertext
			key:           "shortkey",                                     // Invalid key (not 16 bytes)
			expectedPlain: "",
			expectError:   true,
		},
		{
			name:          "Invalid base64 ciphertext",
			ciphertext:    "invalidbase64",    // Invalid base64
			key:           "1234567890123456", // Valid 16-byte key
			expectedPlain: "",
			expectError:   true,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the DecryptAES128ECB function
			plaintext, err := DecryptAES128ECB(tt.ciphertext, tt.key)

			// Check for errors
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, plaintext)
			} else {
				assert.NoError(t, err)
				// Check if the plaintext matches the expected value
				assert.Equal(t, tt.expectedPlain, string(plaintext))
			}
		})
	}
}
