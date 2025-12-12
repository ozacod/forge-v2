package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ozacod/cpx/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveTargetConfig(t *testing.T) {
	tests := []struct {
		name           string
		targetName     string
		expectedTarget config.CITarget
	}{
		{
			name:       "Linux AMD64",
			targetName: "linux-amd64",
			expectedTarget: config.CITarget{
				Name:       "linux-amd64",
				Dockerfile: "Dockerfile.linux-amd64",
				Image:      "cpx-linux-amd64",
				Platform:   "linux/amd64",
			},
		},
		{
			name:       "Linux ARM64",
			targetName: "linux-arm64",
			expectedTarget: config.CITarget{
				Name:       "linux-arm64",
				Dockerfile: "Dockerfile.linux-arm64",
				Image:      "cpx-linux-arm64",
				Platform:   "linux/arm64",
			},
		},
		{
			name:       "Linux AMD64 MUSL",
			targetName: "linux-amd64-musl",
			expectedTarget: config.CITarget{
				Name:       "linux-amd64-musl",
				Dockerfile: "Dockerfile.linux-amd64-musl",
				Image:      "cpx-linux-amd64-musl",
				Platform:   "linux/amd64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deriveTargetConfig(tt.targetName)
			assert.Equal(t, tt.expectedTarget.Name, result.Name)
			assert.Equal(t, tt.expectedTarget.Dockerfile, result.Dockerfile)
			assert.Equal(t, tt.expectedTarget.Image, result.Image)
			assert.Equal(t, tt.expectedTarget.Platform, result.Platform)
		})
	}
}

func TestSaveCIConfig(t *testing.T) {
	tmpDir := t.TempDir()
	ciPath := filepath.Join(tmpDir, "cpx.ci")

	// Create test config
	ciConfig := &config.CIConfig{
		Targets: []config.CITarget{
			{
				Name:       "linux-amd64",
				Dockerfile: "Dockerfile.linux-amd64",
				Image:      "cpx-linux-amd64",
				Triplet:    "x64-linux",
				Platform:   "linux/amd64",
			},
		},
		Build: config.CIBuild{
			Type:         "Release",
			Optimization: "2",
			Jobs:         0,
		},
		Output: ".bin/ci",
	}

	// Save config
	err := config.SaveCI(ciConfig, ciPath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(ciPath)
	require.NoError(t, err)

	// Load it back
	loadedConfig, err := config.LoadCI(ciPath)
	require.NoError(t, err)

	// Verify content
	assert.Len(t, loadedConfig.Targets, 1)
	assert.Equal(t, "linux-amd64", loadedConfig.Targets[0].Name)
	assert.Equal(t, "Release", loadedConfig.Build.Type)
	assert.Equal(t, ".bin/ci", loadedConfig.Output)
}

func TestRunAddTargetWithArgs(t *testing.T) {
	// Setup: create temp dir with mock dockerfiles
	tmpDir := t.TempDir()
	dockerfilesDir := filepath.Join(tmpDir, ".config", "cpx", "dockerfiles")
	require.NoError(t, os.MkdirAll(dockerfilesDir, 0755))

	// Create mock Dockerfiles
	require.NoError(t, os.WriteFile(filepath.Join(dockerfilesDir, "Dockerfile.linux-arm64"), []byte("FROM ubuntu"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dockerfilesDir, "Dockerfile.linux-amd64"), []byte("FROM ubuntu"), 0644))

	// Change HOME to temp dir for test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Change to temp dir for cpx.ci output
	projectDir := filepath.Join(tmpDir, "project")
	require.NoError(t, os.MkdirAll(projectDir, 0755))
	oldWd, _ := os.Getwd()
	os.Chdir(projectDir)
	defer os.Chdir(oldWd)

	// Test: add linux-arm64 target via args
	err := runAddTarget(nil, []string{"linux-arm64"})
	require.NoError(t, err)

	// Verify cpx.ci was created with correct target
	ciConfig, err := config.LoadCI("cpx.ci")
	require.NoError(t, err)
	require.Len(t, ciConfig.Targets, 1)
	assert.Equal(t, "linux-arm64", ciConfig.Targets[0].Name)
	assert.Equal(t, "linux/arm64", ciConfig.Targets[0].Platform)

	// Test: add another target
	err = runAddTarget(nil, []string{"linux-amd64"})
	require.NoError(t, err)

	// Verify both targets exist
	ciConfig, err = config.LoadCI("cpx.ci")
	require.NoError(t, err)
	require.Len(t, ciConfig.Targets, 2)

	// Test: adding duplicate should skip
	err = runAddTarget(nil, []string{"linux-arm64"})
	require.NoError(t, err)

	// Should still have 2 targets (not 3)
	ciConfig, err = config.LoadCI("cpx.ci")
	require.NoError(t, err)
	assert.Len(t, ciConfig.Targets, 2)

	// Test: invalid target should error
	err = runAddTarget(nil, []string{"invalid-target"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown target")
}
