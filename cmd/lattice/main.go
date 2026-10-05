package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zro404/lattice/internal/config"
	"github.com/zro404/lattice/internal/logger"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
	Commit    = "none"
)

func main() {
	version := flag.Bool("version", false, "show version")
	flag.Parse()

	if *version {
		fmt.Printf(
			"Version: %s\nCommit: %s\nBuild Date: %s\n",
			Version,
			Commit,
			BuildDate,
		)
		return
	}

	logger.Init()

	logger.Printf("Version: %s, Commit: %s, BuildDate: %s", Version, Commit, BuildDate)

	cfg, err := config.LoadFile("lattice.yaml")
	if err != nil {
		logger.Fatalf("%s", err.Error())
	}

	logger.Printf("Config: %+v", cfg)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	ticker := time.NewTicker(cfg.TickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Printf("Tick")

		case <-ctx.Done():
			logger.Printf("Shutting down...")
			return
		}

	}

}
