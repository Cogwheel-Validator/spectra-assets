package config

type ChainConfig struct {
	PrettyName string `json:"pretty_name"`
	// ChainPath refers to The Spectra explorer path
	ChainPath string   `json:"chain_path"`
	ChainId   string   `json:"chain_id"`
	ChainType string   `json:"chain_type"`
	APIs      []string `json:"apis"`
	RPCs      []string `json:"rpcs"`
}
