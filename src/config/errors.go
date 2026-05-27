package config

import "errors"

var (
	ErrConfigNotFound        = errors.New("config not found")
	ErrConfigDirCannotRead   = errors.New("config directory cannot be read")
	ErrConfigUnsupportedType = errors.New("unsupported chain type")
	ErrConfigLoadingJson     = errors.New("failed to load json config")
	ErrConfigInvalidJson     = errors.New("invalid json config")
)
