package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yechentide/necrack/netease"
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
		worldDir := args[0]

		if _, err := os.Stat(worldDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: World directory '%s' does not exist\n", worldDir)
			os.Exit(1)
		}

		fmt.Printf("Decrypting NetEase world: %s\n", worldDir)

		if err := netease.DecryptWorldDB(worldDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Decryption completed successfully!")
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
