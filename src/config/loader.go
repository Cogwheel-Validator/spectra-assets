package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func LoadConfigs(path string) ([]ChainConfig, error) {
	var configs = make([]ChainConfig, 0)
	loadedDir, err := os.ReadDir(path)

	if err != nil {
		logger.Error("failed to read config directory", "path", path, "error", err)
		return nil, fmt.Errorf("failed to read config directory %s, error: %w", path, err)
	}

	for _, file := range loadedDir {

		if file.IsDir() {
			subConfigs, err := LoadConfigs(fmt.Sprintf("%s/%s", path, file.Name()))
			if err != nil {
				return nil, err
			}
			configs = append(configs, subConfigs...)
		} else if strings.HasSuffix(file.Name(), ".json") {
			jsonFile, err := os.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))

			if err != nil {
				logger.Error("failed to read json config", "path", fmt.Sprintf("%s/%s", path, file.Name()), "error", err)
				return nil, fmt.Errorf("failed to read json config %s, error: %w", file.Name(), err)
			}

			if err := unmarshalJson(path, jsonFile, file, &configs); err != nil {
				return nil, err
			}

		} else {
			logger.Warn("skipped", "path", fmt.Sprintf("%s/%s", path, file.Name()), "reason", "not a json file", "fileType", file.Type())
		}

	}

	return configs, nil
}

func unmarshalJson(
	path string,
	data []byte,
	file os.DirEntry,
	configs *[]ChainConfig,
) error {
	var chainConfig ChainConfig

	if err := json.Unmarshal(data, &chainConfig); err != nil {
		logger.Error("failed to unmarshal json config", "path", fmt.Sprintf("%s/%s", path, file.Name()), "error", err)
		return fmt.Errorf("failed to unmarshal json config %s, error: %w", file.Name(), err)
	}

	if err := validateConfig(&chainConfig, configs); err != nil {
		logger.Error("skipped", "chain_id", chainConfig.ChainId, "chain_type", chainConfig.ChainType, "reason", err)
		return err
	}

	logger.Info("loaded", "chain_id", chainConfig.ChainId, "chain_type", chainConfig.ChainType)
	return nil
}

func validateConfig(config *ChainConfig, configs *[]ChainConfig) error {
	conditions := []func() bool{
		func() bool { return config.ChainType == "bft" || config.ChainType == "tm2" },
		func() bool { return len(config.APIs) > 0 || len(config.RPCs) > 0 },
		func() bool {
			return len(config.PrettyName) > 0 &&
				len(config.ChainId) > 0 &&
				len(config.ChainPath) > 0
		},
	}
	for _, condition := range conditions {
		if !condition() {
			return fmt.Errorf("chain_type is not bft or tm2, chain_id: %s, chain_type: %s, %w",
				config.ChainId, config.ChainType, ErrConfigUnsupportedType,
			)
		}
	}
	*configs = append(*configs, *config)
	return nil
}
