package server

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yechentide/necrack/netease"
)

func DecryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("zipfile")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		http.Error(w, "File must be a ZIP archive", http.StatusBadRequest)
		return
	}

	tempDir, err := os.MkdirTemp("", "necrack-*")
	if err != nil {
		http.Error(w, "Failed to create temp directory", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	tempZipPath := filepath.Join(tempDir, "input.zip")
	tempZipFile, err := os.Create(tempZipPath)
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tempZipFile, file)
	tempZipFile.Close()
	if err != nil {
		http.Error(w, "Failed to save uploaded file", http.StatusInternalServerError)
		return
	}

	extractDir := filepath.Join(tempDir, "extracted")
	if err := extractZip(tempZipPath, extractDir); err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract ZIP: %v", err), http.StatusInternalServerError)
		return
	}

	worldDirs, err := findWorldDirectories(extractDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find world directories: %v", err), http.StatusInternalServerError)
		return
	}

	if len(worldDirs) == 0 {
		http.Error(w, "No world directories found in ZIP", http.StatusBadRequest)
		return
	}

	for _, worldDir := range worldDirs {
		if err := netease.DecryptWorldDB(worldDir); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decrypt world: %v", err), http.StatusInternalServerError)
			return
		}
	}

	outputZipPath := filepath.Join(tempDir, "decrypted.zip")
	if err := createZip(extractDir, outputZipPath); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create output ZIP: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=decrypted_"+header.Filename)

	outputFile, err := os.Open(outputZipPath)
	if err != nil {
		http.Error(w, "Failed to open output file", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(w, outputFile)
	if err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			rc.Close()
			return fmt.Errorf("invalid file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.FileInfo().Mode())
			rc.Close()
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			rc.Close()
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func createZip(src, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func findWorldDirectories(root string) ([]string, error) {
	var worldDirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "db" {
			worldDir := filepath.Dir(path)
			worldDirs = append(worldDirs, worldDir)
		}

		return nil
	})

	return worldDirs, err
}
