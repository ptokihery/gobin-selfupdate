package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputPath  string
	outputPath string
	keyHex     string
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a config.json file to config.json.enc using AES-256",
	Long: `Encrypt the input config.json using a 256-bit AES key (provided in hex format),
and output the result as an encrypted config.json.enc file.`,
	Example: `
# Basic usage
encryptor encrypt -i config.json -o config.json.enc -k 1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
`,
	Run: func(cmd *cobra.Command, args []string) {
		runEncrypt()
	},
}

func init() {
	encryptCmd.Flags().StringVarP(&inputPath, "in", "i", "", "Input config.json file (required)")
	encryptCmd.Flags().StringVarP(&outputPath, "out", "o", "", "Output encrypted file (required)")
	encryptCmd.Flags().StringVarP(&keyHex, "key", "k", "", "AES-256 encryption key in hex (required)")

	encryptCmd.MarkFlagRequired("in")
	encryptCmd.MarkFlagRequired("out")
	encryptCmd.MarkFlagRequired("key")
}

func runEncrypt() {
	key, err := hex.DecodeString(keyHex)
	if err != nil || len(key) != 32 {
		fmt.Fprintln(os.Stderr, "❌ Invalid key: must be 64 hex characters (32 bytes)")
		os.Exit(1)
	}

	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to read input: %v\n", err)
		os.Exit(1)
	}

	ciphertext, err := encryptAES(plaintext, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Encryption error: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, ciphertext, 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to write encrypted file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Config encrypted successfully:", outputPath)
}

func encryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}
