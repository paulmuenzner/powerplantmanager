package crypto

import (
	"crypto/rand"
	"encoding/hex"
)

type AllowedByteSize int

const (
	ByteSize18  AllowedByteSize = 18
	ByteSize32  AllowedByteSize = 32
	ByteSize64  AllowedByteSize = 64
	ByteSize128 AllowedByteSize = 128
	ByteSize256 AllowedByteSize = 256
)

// GenerateKey generates a key with the specified byte size
func (bs AllowedByteSize) GenerateKey() (string, error) {
	byteKey := make([]byte, int64(bs))

	_, err := rand.Read(byteKey)
	if err != nil {
		return "", err
	}
	hexString := hex.EncodeToString(byteKey)
	return hexString, nil
}
