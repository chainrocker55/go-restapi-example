package utils

import (
	"bytes"
	"creditlimit-connector/app/configs"
	"creditlimit-connector/app/models"
	"errors"
	"fmt"

	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"time"

	"creditlimit-connector/app/log"
)

func EncryptSBAAuthKey(config configs.SBA) string {
	plaintext := GenerateSBAAuthKeyFormattedString(config)
	encText, err := EncryptAES128ECB([]byte(plaintext), config.EncryptionKey)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return encText
}

func EncryptSBARequest[T any](req T, config configs.SBA) string {
	// Convert (marshal) the struct to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	encText, err := EncryptAES128ECB(reqJSON, config.MessageKey)
	if err != nil {
		log.Error(err)
		panic(models.NewErrorResponse(500, "500", "Cannot EncryptSBARequest"))
	}
	return encText
}

func DecryptSBARequest(req models.EncryptQueryCreditLimit, config configs.SBA) []byte {
	// Convert (marshal) the struct to JSON
	decrypted, err := DecryptAES128ECB(req.Msg, config.MessageKey)
	if err != nil {
		log.Error(err)
		panic(models.NewErrorResponse(500, "500", "Cannot DecryptSBARequest"))
	}
	return decrypted
}

func GenerateSBAAuthKeyFormattedString(config configs.SBA) string {
	// Get the current date and time in "YYYYMMDD HH:mm:ss" format.
	serverDate := time.Now().Format("20060102 15:04:05")

	// Generate random values for {random1}, {random2}, {random3}
	random1 := rand.Intn(100) // Generates a random number between 0 and 99
	random2 := rand.Intn(100)
	random3 := rand.Intn(100)

	// Retrieve environment variables (assuming they are set in the system)
	SBA_PID := config.PID
	SBA_HASHWORD := config.Hashword
	SBA_MACADDRESS := config.MACAddress

	// Check if any required environment variables are missing
	if SBA_PID == "" || SBA_HASHWORD == "" || SBA_MACADDRESS == "" {
		return "Error: Missing environment variables"
	}

	// Format the string using the gathered values
	return fmt.Sprintf(
		"a:%d||PID:%s||a:%d||Hashword:%s||a:%d||MacAddress:%s||a:%d||ServerDate:%s",
		random1, SBA_PID, random2, SBA_HASHWORD, random3, SBA_MACADDRESS, random3, serverDate,
	)
}

// NewPkcs7Padding adds PKCS#7 padding to the plaintext to make it a multiple of block size (16 bytes)
func NewPkcs7Padding(plaintext []byte, blockSize int) []byte {
	// Calculate the number of padding bytes needed
	padding := blockSize - len(plaintext)%blockSize
	if padding == 0 {
		padding = blockSize
	}

	// Create the padding
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	// Return the padded plaintext
	return append(plaintext, padText...)
}

// EncryptAES128ECB encrypts plaintext using AES-128 in ECB mode.
func EncryptAES128ECB(plaintext []byte, key string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 {
		return "", errors.New("AES key must be 16 bytes long for AES-128")
	}

	// Pad plaintext to ensure it's a multiple of block size (16 bytes for AES)

	// Create the AES cipher block using the key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	paddedPlaintext := NewPkcs7Padding(plaintext, block.BlockSize())
	// Encrypt the padded plaintext in ECB mode
	ciphertext := make([]byte, len(paddedPlaintext))
	for i := 0; i < len(paddedPlaintext); i += block.BlockSize() {
		block.Encrypt(ciphertext[i:i+block.BlockSize()], paddedPlaintext[i:i+block.BlockSize()])
	}

	// Return the ciphertext encoded in base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES128ECB decrypts ciphertext using AES-128 in ECB mode.
func DecryptAES128ECB(ciphertextBase64, key string) ([]byte, error) {
	// Convert the key to a 16-byte slice (128 bits)
	keyBytes := []byte(key)
	if len(keyBytes) != 16 {
		return []byte(""), fmt.Errorf("key must be 16 bytes (128 bits) long, got %d bytes", len(keyBytes))
	}

	// Decode the base64-encoded ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return []byte(""), err
	}

	// Create a new AES cipher block using the key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return []byte(""), err
	}

	blockSize := block.BlockSize()
	// Decrypt the ciphertext in ECB mode
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += blockSize {
		block.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
	}
	padding := plaintext[len(plaintext)-1]
	return plaintext[:len(plaintext)-int(padding)], nil
}
