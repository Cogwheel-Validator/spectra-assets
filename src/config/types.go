package config

import "os"

type ChainConfig struct {
	PrettyName string `json:"pretty_name"`
	// ChainPath refers to The Spectra explorer path
	ChainPath string   `json:"chain_path"`
	ChainId   string   `json:"chain_id"`
	ChainType string   `json:"chain_type"`
	APIs      []string `json:"apis"`
	RPCs      []string `json:"rpcs"`
	// DataDir is the directory where this chain's data lives; set by the loader.
	DataDir string `json:"-"`
}

type ValidatorData struct {
	Address         string
	Moniker         *string `json:"moniker,omitempty"`
	Identity        *string `json:"identity,omitempty"`
	Description     *string `json:"description,omitempty"`
	Website         *string `json:"website,omitempty"`
	SecurityContact *string `json:"security_contact,omitempty"`
	Logo            *os.File
}
