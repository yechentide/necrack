package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/yechentide/necrack/netease"
	"github.com/yechentide/necrack/styles"
)

var encodeCmd = &cobra.Command{
	Use:   "encode [file|world directory] [key]",
	Short: "Encrypt files using NetEase format",
	Long: `Encrypt files using NetEase Minecraft's custom encryption format.

This command takes a file or world directory and encrypts it with NetEase's
encryption algorithm, making it compatible with NetEase Minecraft world
database format.

If netease.key exists in the world directory, next to the file, or in the
parent world directory when encrypting a file under db, that key is used.
Otherwise, the key argument is used. If neither exists, the default key is used.

Example:
  necrack encode leveldb_file.ldb 1a2b3c4d5e6f7a8b
  necrack encode ./world/db/000005.ldb
  necrack encode ./world`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		filePath := args[0]
		keyHex := ""
		if len(args) == 2 {
			keyHex = args[1]
		}

		// Setup styled output from centralized styles

		// Setup logger
		logger := log.NewWithOptions(nil, log.Options{
			ReportTimestamp: true,
			TimeFormat:      "15:04:05",
			Prefix:          "[encode]",
		})

		logger.Info("Starting file encryption", "file_path", filePath, "key_length", len(keyHex))

		fmt.Println(styles.EncodeHeaderStyle.Render("🔒 NetEase File Encryption"))
		fmt.Printf("File: %s\n", styles.PathStyle.Render(filePath))

		fileInfo, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			logger.Error("File does not exist", "file_path", filePath)
			fmt.Fprintf(os.Stderr, "❌ Error: File '%s' does not exist\n", filePath)
			os.Exit(1)
		} else if err != nil {
			logger.Error("Failed to inspect path", "file_path", filePath, "error", err)
			fmt.Fprintf(os.Stderr, "❌ Error: Failed to inspect path '%s': %v\n", filePath, err)
			os.Exit(1)
		}

		var key []byte
		if keyPath, ok := netease.FindKeyFileForPath(filePath); ok {
			var err error
			key, err = netease.LoadKeyFile(keyPath)
			if err != nil {
				logger.Error("Invalid key file", "key_path", keyPath, "error", err)
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}

			logger.Info("Loaded key from file", "key_path", keyPath)
			fmt.Printf("Key:  %s\n\n", styles.PathStyle.Render(keyPath))
		} else {
			if keyHex == "" {
				keyHex = netease.DefaultKeyHex
				logger.Info("Using default key", "key_file", netease.KeyFileName)
			} else {
				logger.Info("Using key from argument")
			}

			var err error
			key, err = netease.ParseHexKey(keyHex)
			if err != nil {
				logger.Error("Invalid key format", "key_hex", keyHex, "error", err)
				fmt.Fprintf(os.Stderr, "❌ Error: Invalid key format: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Key:  %s\n\n", styles.KeyStyle.Render(keyHex))
		}

		logger.Info("Key parsed successfully, starting encryption")

		if fileInfo.IsDir() {
			encryptedDir, err := netease.EncryptWorldDB(filePath, key)
			if err != nil {
				logger.Error("World encryption failed", "world_dir", filePath, "error", err)
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}

			duration := time.Since(start)
			logger.Info("World encryption completed successfully",
				"input_dir", filePath,
				"output_dir", encryptedDir,
				"duration", duration)

			fmt.Println(styles.SuccessStyle.Render("✅ Encryption completed successfully!"))
			fmt.Printf("📁 Encrypted world saved to: %s\n", styles.PathStyle.Render(encryptedDir))
			fmt.Printf("⏱️  Completed in %v\n", duration)
			return
		}

		encrypted, err := netease.EncryptFile(filePath, key)
		if err != nil {
			logger.Error("Encryption failed", "file_path", filePath, "error", err)
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		outputPath := filePath + ".encrypted"
		if err := os.WriteFile(outputPath, encrypted, 0644); err != nil {
			logger.Error("Failed to write encrypted file", "output_path", outputPath, "error", err)
			fmt.Fprintf(os.Stderr, "❌ Error writing encrypted file: %v\n", err)
			os.Exit(1)
		}

		duration := time.Since(start)
		logger.Info("Encryption completed successfully",
			"input_file", filePath,
			"output_file", outputPath,
			"file_size", len(encrypted),
			"duration", duration)

		fmt.Println(styles.SuccessStyle.Render("✅ Encryption completed successfully!"))
		fmt.Printf("📄 Output: %s\n", styles.PathStyle.Render(outputPath))
		fmt.Printf("⏱️  Completed in %v\n", duration)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
