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

var decodeCmd = &cobra.Command{
	Use:   "decode [world directory]",
	Short: "Decrypt NetEase Minecraft world files",
	Long: `Decrypt NetEase Minecraft world files in the specified world directory.
The world directory should contain a 'db' subdirectory with encrypted files.

Example:
  necrack decode ./ne-worlds/661428f7-1e29-47ca-99af-c1eac0c41ba5`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		worldDir := args[0]
		
		// Setup styled output from centralized styles
		
		// Setup logger
		logger := log.NewWithOptions(nil, log.Options{
			ReportTimestamp: true,
			TimeFormat:      "15:04:05",
			Prefix:          "[decode]",
		})
		
		logger.Info("Starting world decryption", "world_dir", worldDir)
		
		fmt.Println(styles.DecodeHeaderStyle.Render("🔓 NetEase World Decryption"))
		fmt.Printf("Target: %s\n\n", styles.PathStyle.Render(worldDir))

		if _, err := os.Stat(worldDir); os.IsNotExist(err) {
			logger.Error("World directory does not exist", "world_dir", worldDir)
			fmt.Fprintf(os.Stderr, "❌ Error: World directory '%s' does not exist\n", worldDir)
			os.Exit(1)
		}

		logger.Info("World directory found, starting decryption process")

		if err := netease.DecryptWorldDB(worldDir); err != nil {
			logger.Error("Decryption failed", "world_dir", worldDir, "error", err)
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

		duration := time.Since(start)
		logger.Info("Decryption completed successfully", "world_dir", worldDir, "duration", duration)
		fmt.Println(styles.SuccessStyle.Render("✅ Decryption completed successfully!"))
		fmt.Printf("⏱️  Completed in %v\n", duration)
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
