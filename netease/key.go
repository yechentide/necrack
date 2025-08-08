package netease

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DeriveKey(dbDir string) ([]byte, error) {
	manifestName, err := findManifestFile(dbDir)
	if err != nil {
		return nil, fmt.Errorf("failed to find MANIFEST file: %w", err)
	}

	currentPath := filepath.Join(dbDir, "CURRENT")
	currentData, err := os.ReadFile(currentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CURRENT file: %w", err)
	}

	if err := ValidateDecryptableFile(currentData); err != nil {
		return nil, fmt.Errorf("CURRENT file is not decryptable: %w", err)
	}

	currentBody := currentData[4:]

	// Add newline to manifestName
	manifestWithNewline := append(manifestName, '\n')

	if len(manifestWithNewline) != 16 {
		return nil, fmt.Errorf("manifest name with newline must be exactly 16 bytes, got %d", len(manifestWithNewline))
	}

	if len(currentBody) < 16 {
		return nil, fmt.Errorf("current body must be at least 16 bytes, got %d", len(currentBody))
	}

	keyRaw := make([]byte, 16)
	for i := 0; i < 16; i++ {
		keyRaw[i] = currentBody[i] ^ manifestWithNewline[i]
	}

	first8 := keyRaw[:8]
	last8 := keyRaw[8:]

	if !bytes.Equal(first8, last8) {
		return nil, fmt.Errorf("key verification failed: first 8 bytes do not match last 8 bytes")
	}

	return first8, nil
}

func findManifestFile(dbDir string) ([]byte, error) {
	entries, err := os.ReadDir(dbDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read db directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), "MANIFEST-") {
			return []byte(entry.Name()), nil
		}
	}

	return nil, fmt.Errorf("no MANIFEST file found in db directory")
}
