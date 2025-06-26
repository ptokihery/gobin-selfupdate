package main

import (
	"log"

	"github.com/ptokihery/gobin-selfupdate/cmd/encryptor/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

