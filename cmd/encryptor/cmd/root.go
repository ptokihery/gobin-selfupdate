package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "encryptor",
	Short: "A CLI tool to encrypt config files using AES-256",
	Long: `Encryptor is a simple CLI tool to protect your config.json
by encrypting it with a 256-bit AES key. The result can be embedded
in Go applications securely and decrypted at runtime.`,
	Example: `
# Encrypt config.json using a 32-byte AES key (hex format)
encryptor encrypt -i config.json -o config.json.enc -k 1234567890abcdef...`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}
