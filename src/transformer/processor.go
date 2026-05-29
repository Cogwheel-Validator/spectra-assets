package transformer

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"spectra-assets/src/config"
	"spectra-assets/src/fetch"
)

// ProcessAll processes every chain in configs and writes assets under assetsPath.
func ProcessAll(assetsPath string, configs []config.ChainConfig) error {
	for _, cfg := range configs {
		if err := processChain(assetsPath, cfg); err != nil {
			return fmt.Errorf("chain %s: %w", cfg.ChainId, err)
		}
	}
	return nil
}

func processChain(assetsPath string, cfg config.ChainConfig) error {
	entries, err := LoadValidators(cfg)
	if err != nil {
		return fmt.Errorf("loading validators: %w", err)
	}

	chainDir := filepath.Join(assetsPath, cfg.ChainPath)
	if err := os.MkdirAll(chainDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir %s: %w", chainDir, err)
	}

	for _, entry := range entries {
		if err := processValidator(entry, chainDir); err != nil {
			logger.Warn("skipping validator", "address", entry.Address, "chain", cfg.ChainId, "error", err)
		}
	}
	return nil
}

func processValidator(entry ValidatorEntry, assetsDir string) error {
	asset := ValidatorAsset{
		Address:     entry.Address,
		Moniker:     entry.Meta.Moniker,
		Identity:    entry.Meta.Identity,
		Description: entry.Meta.Description,
	}

	jsonData, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	jsonPath := filepath.Join(assetsDir, entry.Address+".json")
	if err := os.WriteFile(jsonPath, jsonData, 0o644); err != nil {
		return fmt.Errorf("write JSON: %w", err)
	}
	logger.Info("wrote", "file", jsonPath)

	img, err := resolveImage(entry)
	if err != nil {
		logger.Warn("no logo available", "address", entry.Address, "error", err)
		return nil
	}

	processed := resizeCover(img, logoSize)
	jpegData, err := encodeJPEG(processed)
	if err != nil {
		return fmt.Errorf("encode JPEG: %w", err)
	}

	jpegPath := filepath.Join(assetsDir, entry.Address+".jpg")
	if err := os.WriteFile(jpegPath, jpegData, 0o644); err != nil {
		return fmt.Errorf("write JPEG: %w", err)
	}
	logger.Info("wrote", "file", jpegPath)
	return nil
}

// resolveImage returns the validator logo, trying the local file first then falling back to Keybase.
func resolveImage(entry ValidatorEntry) (image.Image, error) {
	if entry.LogoPath != "" {
		decoded, err := decodeImageFile(entry.LogoPath)
		if err == nil {
			return decoded, nil
		}
		logger.Warn("logo file invalid", "path", entry.LogoPath, "error", err)
	}

	if entry.Meta.Identity == "" {
		return nil, fmt.Errorf("no logo and no Keybase identity")
	}

	resp, err := fetch.QueryKeybase(entry.Meta.Identity)
	if err != nil {
		return nil, fmt.Errorf("Keybase query: %w", err)
	}
	if len(resp.Them) == 0 || resp.Them[0].Pictures.Primary.Url == "" {
		return nil, fmt.Errorf("no picture in Keybase response for identity %s", entry.Meta.Identity)
	}

	decoded, err := fetchImageFromURL(resp.Them[0].Pictures.Primary.Url)
	if err != nil {
		return nil, fmt.Errorf("fetching Keybase image: %w", err)
	}
	return decoded, nil
}
