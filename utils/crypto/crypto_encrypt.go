package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// EncryptWithKey encrypts plaintext using the provided encryption key and returns the ciphertext and IV.
// It uses AES-GCM mode for encryption.
func EncryptWithKey(plaintext string, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv, err := GenerateIV(aes.BlockSize)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, iv, []byte(plaintext), nil)
	return ciphertext, iv, nil
}

// Example usage:
// ciphertext, iv, err := crypto.EncryptWithKey(plaintext, key)
// if err != nil {
//     panic(err)
// }
// fmt.Printf("Ciphertext: %v\n", ciphertext)
// fmt.Printf("IV: %v\n", iv)

// GenerateIV generates a random Initialization Vector (IV) with the specified size.
func GenerateIV(size int) ([]byte, error) {
	iv := make([]byte, size)
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	return iv, nil
}

// Example with AWS KSM key
// keyID := "your_key_id_here"

// // Retrieve the encryption key
// keyBytes, err := RetrieveEncryptionKey(keyID)
// if err != nil {
//     log.Fatal(err)
// }

// // Encrypt plaintext using the retrieved key
// plaintext := "your_plaintext_here"
// ciphertext, iv, err := EncryptWithKey(plaintext, keyBytes)
// if err != nil {
//     log.Fatal(err)
// }

// Now you have the ciphertext and IV for further processing.
