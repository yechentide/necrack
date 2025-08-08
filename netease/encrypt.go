package netease

import (
	"encoding/hex"
	"fmt"
	"os"
)

func EncryptFile(filePath string, key []byte) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	result := encryptData(data, key)
	return result, nil
}

func encryptData(data []byte, key []byte) []byte {
	encrypted := xorDecrypt(data, key)

	result := make([]byte, 0, len(headerNetEaseCurrent)+len(encrypted))
	result = append(result, headerNetEaseCurrent...)
	result = append(result, encrypted...)

	return result
}

func ParseHexKey(keyHex string) ([]byte, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	if len(key) == 0 {
		return nil, fmt.Errorf("key cannot be empty")
	}

	if len(key) != 8 {
		return nil, fmt.Errorf("key must be exactly 8 bytes (16 hex characters), got %d bytes", len(key))
	}

	return key, nil
}
