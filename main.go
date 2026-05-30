package main

import (
	"log/slog"
	"os"
	"spectra-assets/src/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}
