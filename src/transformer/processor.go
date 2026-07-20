package transformer

import (
	"encoding/json"
	"fmt"
	"image"
	"maps"
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
	chainDir := filepath.Join(assetsPath, cfg.ChainPath)
	if err := os.MkdirAll(chainDir, 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}
	if cfg.ChainType == "bft" {
		return processChainBFT(chainDir, cfg)
	}
	return processChainDataDir(chainDir, cfg)
}

// processChainBFT fetches all on-chain validators, then merges any data-dir overrides on top.
// Data-dir entries win for every field they define; on-chain values fill gaps.
func processChainBFT(chainDir string, cfg config.ChainConfig) error {
	// Collect data-dir entries (these are PR-submitted overrides/additions).
	dataEntries := map[string]ValidatorEntry{}
	if entries, err := LoadValidators(cfg); err == nil {
		for _, e := range entries {
			dataEntries[e.Address] = e
		}
	}

	fetcher, err := fetch.NewCosmosFetcher(cfg.ChainId, cfg.APIs)
	if err != nil {
		logger.Warn("chain API unavailable, falling back to data-dir only", "chain", cfg.ChainId, "error", err)
		for _, entry := range dataEntries {
			processValidatorLogged(entry, chainDir, cfg.ChainId)
		}
		return nil
	}

	chainVals := fetcher.GatherChainOperators(false)
	if cfg.Governors {
		// Should be okay since we will just mix up valopers and gov addresses.
		chainGovs := fetcher.GatherChainOperators(true)
		maps.Copy(chainVals, chainGovs)
	}

	processed := make(map[string]bool, len(chainVals))

	for address, desc := range chainVals {
		entry := ValidatorEntry{
			Address: address,
			Meta: ValidatorMeta{
				Moniker:         desc.Moniker,
				Identity:        desc.Identity,
				Description:     desc.Details,
				Website:         desc.Website,
				SecurityContact: desc.SecurityContact,
			},
		}
		// Data-dir entry overrides on-chain fields where it provides values.
		if override, ok := dataEntries[address]; ok {
			if override.Meta.Moniker != "" {
				entry.Meta.Moniker = override.Meta.Moniker
			}
			if override.Meta.Identity != "" {
				entry.Meta.Identity = override.Meta.Identity
			}
			if override.Meta.Description != "" {
				entry.Meta.Description = override.Meta.Description
			}
			if override.Meta.Website != "" {
				entry.Meta.Website = override.Meta.Website
			}
			if override.Meta.SecurityContact != "" {
				entry.Meta.SecurityContact = override.Meta.SecurityContact
			}
			entry.DataDir = override.DataDir
			entry.LogoPath = override.LogoPath
		}
		processValidatorLogged(entry, chainDir, cfg.ChainId)
		processed[address] = true
	}

	// Handle data-dir entries that have no on-chain counterpart (e.g. manual additions).
	for address, entry := range dataEntries {
		if !processed[address] {
			processValidatorLogged(entry, chainDir, cfg.ChainId)
		}
	}
	return nil
}

// processChainDataDir processes TM2 (or any non-BFT) chain using only the data directory.
func processChainDataDir(chainDir string, cfg config.ChainConfig) error {
	entries, err := LoadValidators(cfg)
	if err != nil {
		return fmt.Errorf("loading validators: %w", err)
	}
	for _, entry := range entries {
		processValidatorLogged(entry, chainDir, cfg.ChainId)
	}
	return nil
}

// ProcessEntry processes exactly one data-dir entry for the given chain + address.
// No chain API is called — this is the path used in PR workflows.
func ProcessEntry(assetsPath string, cfg config.ChainConfig, address string) error {
	valDir := filepath.Join(cfg.DataDir, address)
	if _, err := os.Stat(valDir); err != nil {
		return fmt.Errorf("validator directory not found: %s", valDir)
	}

	entry := ValidatorEntry{
		Address: address,
		DataDir: valDir,
	}

	if data, err := os.ReadFile(filepath.Join(valDir, "validator.json")); err == nil {
		var meta ValidatorMeta
		if err := json.Unmarshal(data, &meta); err == nil {
			entry.Meta = meta
		}
	}

	for _, name := range []string{"logo.jpg", "logo.jpeg"} {
		candidate := filepath.Join(valDir, name)
		if _, err := os.Stat(candidate); err == nil {
			entry.LogoPath = candidate
			break
		}
	}

	chainDir := filepath.Join(assetsPath, cfg.ChainPath)
	if err := os.MkdirAll(chainDir, 0o755); err != nil {
		return err
	}
	return processValidator(entry, chainDir)
}

func processValidatorLogged(entry ValidatorEntry, chainDir, chainId string) {
	if err := processValidator(entry, chainDir); err != nil {
		logger.Warn("skipping validator", "address", entry.Address, "chain", chainId, "error", err)
	}
}

func processValidator(entry ValidatorEntry, assetsDir string) error {
	asset := ValidatorAsset{
		Address:         entry.Address,
		Moniker:         entry.Meta.Moniker,
		Identity:        entry.Meta.Identity,
		Description:     entry.Meta.Description,
		Website:         entry.Meta.Website,
		SecurityContact: entry.Meta.SecurityContact,
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

// resolveImage returns the logo image: local file first, then Keybase fallback.
func resolveImage(entry ValidatorEntry) (image.Image, error) {
	if entry.LogoPath != "" {
		if decoded, err := decodeImageFile(entry.LogoPath); err == nil {
			return decoded, nil
		} else {
			logger.Warn("logo file invalid", "path", entry.LogoPath, "error", err)
		}
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
