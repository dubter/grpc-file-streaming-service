package config

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	testConfigPath := "../../config/server.yaml"
	config, err := ParseConfig(context.Background(), testConfigPath)

	require.NoError(t, err)
	require.NotNil(t, config)
}

func TestConfigEmptyPath(t *testing.T) {
	config, err := ParseConfig(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, config)
}
