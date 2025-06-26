package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

func EncryptFile(inputFile string, outputFile string, key []byte) error {
	plaintext, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return os.WriteFile(outputFile, ciphertext, 0644)
}
