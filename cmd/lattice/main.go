package main

import (
	"flag"
	"fmt"

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

	config := config.LoadFile("lattice.yaml")

	logger.Printf("%+v", config)

}
