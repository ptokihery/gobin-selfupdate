# gobin-selfupdate

A professional Go library and CLI tool for securely **self-updating Go applications**. Features include:

- Checking and applying updates from S3 or any HTTP server
- Support for a **manifest file** alongside the binary
- AES-256 **encrypted config files** for secure AWS S3 credentials
- A reusable CLI encryptor tool to generate encrypted config files

---

## Components

### 1. `updater` Package

A Go library to integrate self-update logic into your application.

### 2. `encryptor` CLI

A command-line tool to encrypt your `config.json` into `config.json.enc` using AES-256, ensuring your credentials remain secure.

---

## Installation

### Encryptor CLI

Install the encryptor CLI globally to manage your encrypted configuration files:

```bash
go install github.com/ptokihery/gobin-selfupdate/cmd/encryptor@v0.1.3
```

Make sure your `$GOBIN` is in your `PATH`.

---

## Usage

### 1. Encrypt your config file

Before using the updater, encrypt your plaintext `config.json`:

```bash
encryptor encrypt -i config.json -o config.json.enc -k <32-byte-hex-aes-key>
```

This generates `config.json.enc` which your app will load securely.

---

### 2. Using the updater package in your Go app

Example usage:

```go
package main

import (
    "log"
    "time"

    "github.com/ptokihery/gobin-selfupdate/config"
    "github.com/ptokihery/gobin-selfupdate/updater"
)

func main() {
    // Load encrypted config (must decrypt inside your app)
    cfg, err := config.LoadEncryptedConfig("config.json.enc", "<your-32-byte-hex-aes-key>")
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Create updater instance with config
    u, err := updater.NewUpdater(cfg)
    if err != nil {
        log.Fatalf("failed to create updater: %v", err)
    }

    // Run update check every 10 minutes (customizable interval)
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()

    // Initial check
    go func() {
        if err := u.CheckAndUpdate(); err != nil {
            log.Printf("update check failed: %v", err)
        }
    }()

    for range ticker.C {
        if err := u.CheckAndUpdate(); err != nil {
            log.Printf("update check failed: %v", err)
        }
    }

    // Your app logic here...
}
```

---

## Configuration file format

The config JSON should contain fields like AWS region, credentials, bucket name, object keys, etc.  
This is what you encrypt with the encryptor CLI.

Example `config.json`:

```json
{
  "awsRegion": "us-east-1",
  "accessKey": "AKIAxxxxxxx",
  "secretKey": "xxxxxxxxxx",
  "bucket": "my-bucket",
  "binaryKey": "updates/myapp_latest",
  "manifestKey": "updates/myapp_latest.manifest",
  "outputPath": "/tmp/myapp_new",
  "updateIntervalMinutes": 10
}
```

Encrypt it before deploying.

---

## Why encrypt config?

- Protect your AWS keys and secrets from leaking
- Embed encrypted config in your app binary (via Go embed)
- Decrypt at runtime securely

---

## Author

Tokihery RANDRIANAMBININTSOA, Dev

---

## License

MIT License. See LICENSE for details.