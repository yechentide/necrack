package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "necrack",
	Short: "NetEase Minecraft world file encryption/decryption tool",
	Long: `necrack is a command-line tool for working with NetEase Minecraft world files.

It provides functionality to decrypt and encrypt NetEase Minecraft world database files,
allowing you to work with world data that uses NetEase's custom encryption format.

Available commands:
  decode    Decrypt NetEase Minecraft world files
  encode    Encrypt files using NetEase format

Use "necrack help [command]" for more information about a specific command.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.necrack.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
