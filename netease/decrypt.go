package netease

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
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

func DecryptWorldDB(worldDir string) (string, error) {
	dbDir := filepath.Join(worldDir, "db")
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		return "", fmt.Errorf("db directory not found in %s", worldDir)
	}

	// Create a copy of the world directory
	timestamp := time.Now().Format("20060102_150405")
	worldDirName := filepath.Base(worldDir)
	copyDir := filepath.Join(filepath.Dir(worldDir), worldDirName+"_decrypted_"+timestamp)
	
	if err := copyDirectory(worldDir, copyDir); err != nil {
		return "", fmt.Errorf("failed to copy world directory: %w", err)
	}

	// Work on the copied directory
	copyDbDir := filepath.Join(copyDir, "db")
	key, err := DeriveKey(copyDbDir)
	if err != nil {
		return "", fmt.Errorf("failed to derive key: %w", err)
	}

	err = filepath.WalkDir(copyDbDir, func(path string, d fs.DirEntry, err error) error {
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

		// Overwrite the original file in the copy with decrypted data
		if err := os.WriteFile(path, decrypted, 0644); err != nil {
			return fmt.Errorf("failed to write decrypted file %s: %w", path, err)
		}

		fmt.Printf("Decrypted: %s\n", path)
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to process db directory: %w", err)
	}

	return copyDir, nil
}

func copyDirectory(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", path, err)
		}
		defer srcFile.Close()

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(dstPath), err)
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", dstPath, err)
		}
		defer dstFile.Close()

		srcInfo, err := srcFile.Stat()
		if err != nil {
			return fmt.Errorf("failed to get source file info: %w", err)
		}

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}

		return os.Chmod(dstPath, srcInfo.Mode())
	})
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
