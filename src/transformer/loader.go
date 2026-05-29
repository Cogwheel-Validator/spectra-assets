package transformer

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"spectra-assets/src/config"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// LoadValidators scans a chain's data directory and returns one entry per validator sub-directory.
func LoadValidators(cfg config.ChainConfig) ([]ValidatorEntry, error) {
	dirs, err := os.ReadDir(cfg.DataDir)
	if err != nil {
		return nil, err
	}

	var entries []ValidatorEntry

	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}

		address := d.Name()
		valDir := filepath.Join(cfg.DataDir, address)

		entry := ValidatorEntry{
			Address: address,
			DataDir: valDir,
		}

		if data, err := os.ReadFile(filepath.Join(valDir, "validator.json")); err == nil {
			var meta ValidatorMeta
			if err := json.Unmarshal(data, &meta); err == nil {
				entry.Meta = meta
			} else {
				logger.Warn("invalid validator.json", "address", address, "chain", cfg.ChainId, "error", err)
			}
		}

		for _, name := range []string{"logo.png", "logo.jpg", "logo.jpeg"} {
			candidate := filepath.Join(valDir, name)
			if _, err := os.Stat(candidate); err == nil {
				entry.LogoPath = candidate
				break
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
