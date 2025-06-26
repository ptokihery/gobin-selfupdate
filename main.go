package main
import (
	_ "embed"
	"log"
	"time"

	"github.com/ptokihery/gobin-selfupdate/config"
	"github.com/ptokihery/gobin-selfupdate/updater"
)

//go:embed config.json.enc
var cfgJson []byte

var key = "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef" // AES-256 key

func main() {
	cfg, err := config.Load(cfgJson, key)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	up, err := updater.StartUpdater(cfg, "1.0.0", 10*time.Minute)
	if err != nil {
		log.Fatalf("failed to start updater: %v", err)
	}
	defer up.Stop()

	log.Println("Application is running...")
	log.Printf("CGD : %v", cfg.AWSRegion)
	select {}
}
