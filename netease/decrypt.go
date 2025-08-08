package netease

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func DecryptFile(filePath string, key []byte) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := ValidateDecryptableFile(data); err != nil {
		return nil, fmt.Errorf("file %s is not decryptable: %w", filePath, err)
	}

	body := data[4:]
	decrypted := xorDecrypt(body, key)

	return decrypted, nil
}

func DecryptWorldDB(worldDir string) error {
	dbDir := filepath.Join(worldDir, "db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		return fmt.Errorf("db directory not found in %s", worldDir)
	}

	key, err := DeriveKey(dbDir)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	err = filepath.WalkDir(dbDir, func(path string, d fs.DirEntry, err error) error {
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

		headerType := identifyHeader(data)
		if headerType != HeaderTypeNetEaseCurrent {
			return nil
		}

		decrypted, err := DecryptFile(path, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt file %s: %w", path, err)
		}

		outputPath := path + ".decrypted"
		if err := os.WriteFile(outputPath, decrypted, 0644); err != nil {
			return fmt.Errorf("failed to write decrypted file %s: %w", outputPath, err)
		}

		fmt.Printf("Decrypted: %s -> %s\n", path, outputPath)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to process db directory: %w", err)
	}

	return nil
}

func xorDecrypt(data []byte, key []byte) []byte {
	if len(key) == 0 {
		return data
	}

	result := make([]byte, len(data))
	keyLen := len(key)

	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ key[i%keyLen]
	}

	return result
}
