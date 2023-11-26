package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseConfig(t *testing.T) {
	testConfigPath := "../../config/client.yaml"
	config, err := ParseConfig(testConfigPath)

	require.NoError(t, err)
	require.NotNil(t, config)
}
