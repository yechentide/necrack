package netease

import (
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const DefaultKeyHex = "3838333239383531"

func EncryptFile(filePath string, key []byte) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	result := encryptData(data, key)
	return result, nil
}

func EncryptWorldDB(worldDir string, key []byte) (string, error) {
	dbDir := filepath.Join(worldDir, "db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		return "", fmt.Errorf("db directory not found in %s", worldDir)
	}

	timestamp := time.Now().Format("20060102_150405")
	worldDirName := filepath.Base(worldDir)
	copyDir := filepath.Join(filepath.Dir(worldDir), worldDirName+"_encrypted_"+timestamp)

	if err := copyDirectory(worldDir, copyDir); err != nil {
		return "", fmt.Errorf("failed to copy world directory: %w", err)
	}

	keyPath := filepath.Join(copyDir, KeyFileName)
	if err := os.Remove(keyPath); err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to remove key file from encrypted copy: %w", err)
	}

	copyDbDir := filepath.Join(copyDir, "db")
	err := filepath.WalkDir(copyDbDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if identifyHeader(data) == HeaderTypeNetEaseCurrent {
			return nil
		}

		encrypted := encryptData(data, key)
		if err := os.WriteFile(path, encrypted, 0644); err != nil {
			return fmt.Errorf("failed to write encrypted file %s: %w", path, err)
		}

		fmt.Printf("Encrypted: %s\n", path)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to process db directory: %w", err)
	}

	return copyDir, nil
}

func encryptData(data []byte, key []byte) []byte {
	encrypted := xorDecrypt(data, key)

	result := make([]byte, 0, len(headerNetEaseCurrent)+len(encrypted))
	result = append(result, headerNetEaseCurrent...)
	result = append(result, encrypted...)

	return result
}

func ParseHexKey(keyHex string) ([]byte, error) {
	keyHex = strings.TrimSpace(keyHex)
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

func LoadKeyFile(keyPath string) ([]byte, error) {
	keyHex, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := ParseHexKey(string(keyHex))
	if err != nil {
		return nil, fmt.Errorf("invalid key file %s: %w", keyPath, err)
	}

	return key, nil
}

func FindKeyFileForPath(filePath string) (string, bool) {
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		candidates := []string{
			filepath.Join(filePath, KeyFileName),
			filepath.Join(filepath.Dir(filePath), KeyFileName),
		}

		for _, candidate := range candidates {
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate, true
			}
		}

		return "", false
	}

	fileDir := filepath.Dir(filePath)
	candidates := []string{
		filepath.Join(fileDir, KeyFileName),
	}

	if filepath.Base(fileDir) == "db" {
		candidates = append(candidates, filepath.Join(filepath.Dir(fileDir), KeyFileName))
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}
	}

	return "", false
}
