package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

func plaintexttoken() (string, error) {
	// Generate a random plaintext of length 34
	plaintextBytes := make([]byte, 34)
	_, err := rand.Read(plaintextBytes)
	if err != nil {
		return "", err
	}
	verifyToken := hex.EncodeToString(plaintextBytes)
	return verifyToken, err
}

func GenerateVerifyToken(keyHex string) (string, string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", "", err
	}

	// Check if the key is of the correct size for AES (32 bytes)
	if len(key) != 32 {
		return "", "", errors.New("invalid key size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}
	// Encrypt the plaintext
	plaintext, err := plaintexttoken()
	if err != nil {
		return "", "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode the encrypted ciphertext using base64
	encodedCiphertext := base64.URLEncoding.EncodeToString(ciphertext)

	return plaintext, encodedCiphertext, nil
}

func DecryptVerifyToken(encodedCiphertext, keyHex string) (string, error) {
	// Convert key from hexadecimal to byte slice
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}

	// Check if the key is of the correct size for AES (32 bytes)
	if len(key) != 32 {
		return "", errors.New("invalid key size")
	}

	// Decode the base64-encoded ciphertext
	ciphertext, err := base64.URLEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block using the specified key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM mode with the cipher block
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Check if the ciphertext is valid
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("invalid ciphertext size")
	}

	// Extract the nonce from the ciphertext
	nonce := ciphertext[:gcm.NonceSize()]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
