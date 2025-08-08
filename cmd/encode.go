package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yechentide/necrack/netease"
)

var encodeCmd = &cobra.Command{
	Use:   "encode [file] [key]",
	Short: "Encrypt files using NetEase format",
	Long: `Encrypt files using NetEase Minecraft's custom encryption format.

This command takes a file and encrypts it with NetEase's encryption algorithm,
making it compatible with NetEase Minecraft world database format.

The key should be provided as a hex string (e.g., "1a2b3c4d5e6f7a8b").

Example:
  necrack encode leveldb_file.ldb 1a2b3c4d5e6f7a8b`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		keyHex := args[1]

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", filePath)
			os.Exit(1)
		}

		key, err := netease.ParseHexKey(keyHex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid key format: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Encrypting file: %s\n", filePath)

		encrypted, err := netease.EncryptFile(filePath, key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		outputPath := filePath + ".encrypted"
		if err := os.WriteFile(outputPath, encrypted, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing encrypted file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Encryption completed: %s -> %s\n", filePath, outputPath)
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
}
