package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

// DecryptWithKey decrypts ciphertext using a provided key and IV.
func DecryptWithKey(ciphertext, key, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
