package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zro404/lattice/internal/archive"
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

	handleLogs(cfg)

	// ctx, stop := signal.NotifyContext(
	// 	context.Background(),
	// 	os.Interrupt,
	// 	syscall.SIGTERM,
	// )
	// defer stop()

	// ticker := time.NewTicker(cfg.TickInterval)
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		logger.Printf("Tick")

	// 	case <-ctx.Done():
	// 		logger.Printf("Shutting down...")
	// 		return
	// 	}

	// }

}

func handleLogs(cfg *config.Config) {
	for _, log := range cfg.Logs {
		err := processLog(&log)
		if err != nil {
			logger.Printf("Error processing log %s: %s", log.Name, err.Error())
		}
	}
}

func processLog(log *config.Log) error {
	logger.Printf("Processing log: %s", log.Name)

	tmpFile, err := os.CreateTemp("", fmt.Sprintf("lattice_%s_logs_*.tar.zst", log.Name))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	tmpPath := tmpFile.Name()

	archive, err := archive.NewArchive(tmpFile)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to create archive: %w", err)
	}

	for _, path := range log.Paths {
		logger.Printf("\tPath: %s", path)

		inFile, err := os.Open(path)
		if err != nil {
			archive.Close()
			tmpFile.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("failed to open %s: %w", path, err)
		}

		if err := archive.AddFile(path, inFile); err != nil {
			inFile.Close()
			archive.Close()
			tmpFile.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("failed to add %s to archive: %w", path, err)
		}

		inFile.Close()
	}

	if err := archive.Close(); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close archive: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// TODO: upload tmpPath

	if err := os.Remove(tmpPath); err != nil {
		return fmt.Errorf("failed to remove temporary archive: %w", err)
	}

	logger.Printf("Log backup successful!")

	return nil
}
