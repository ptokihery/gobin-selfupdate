package selfupdate

import (
	"context"
	"log"
	"time"

	"github.com/ptokihery/gobin-seltupdate/internal/updater"
)

type Checker struct {
	runner         *updater.Runner
	currentVersion string
	interval       time.Duration
	cancel         context.CancelFunc
	ctx            context.Context
}

func NewChecker(runner *updater.Runner, currentVersion string, interval time.Duration) *Checker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Checker{
		runner:         runner,
		currentVersion: currentVersion,
		interval:       interval,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (c *Checker) Start() {
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				log.Println("[selfupdate] Update checker stopped")
				return
			case <-ticker.C:
				err := c.runner.Run(c.ctx, c.currentVersion)
				if err != nil {
					log.Printf("[selfupdate] Update check error: %v", err)
				} else {
					log.Println("[selfupdate] Update check completed")
				}
			}
		}
	}()
}

func (c *Checker) Stop() {
	c.cancel()
}
