package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Healthy_Load(t *testing.T) {
	path := "./test_data/tc1"
	require.DirExists(t, path)
	configs, err := LoadConfigs(path)
	if err != nil {
		t.Logf("Error loading test data: %v", err)
	}
	require.Len(t, configs, 2)
}

func Test_Failed_Load(t *testing.T) {
	path := "./test_data/tc2"
	require.DirExists(t, path)
	configs, err := LoadConfigs(path)
	require.Error(t, err)
	require.Nil(t, configs)
}

func Test_Nested_Load(t *testing.T) {
	path := "./test_data/tc3"
	require.DirExists(t, path)
	configs, _ := LoadConfigs(path)
	require.Len(t, configs, 2)
}

func Test_Validation_Fail(t *testing.T) {
	paths := []string{
		"./test_data/tc4/t1",
		"./test_data/tc4/t2",
		"./test_data/tc4/t3",
	}
	for _, path := range paths {
		require.DirExists(t, path)
	}

	var errors = make([]error, 3)
	for i, path := range paths {
		var err error
		config, err := LoadConfigs(path)
		if err != nil {
			errors[i] = err
		}
		require.Nil(t, config)
	}
	require.Len(t, errors, 3)
}
