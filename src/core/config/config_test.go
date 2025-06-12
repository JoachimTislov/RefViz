package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JoachimTislov/RefViz/internal" // For GetAndUnmarshalFile
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetBaseBranchLink(t *testing.T) {
	// Ensure a clean config state for this test
	config = NewConfig()
	originalLink := config.BaseBranchLink // Save original default

	// Test with a valid URL
	validURL := "https://valid.example.com/repo/tree/main/"
	err := SetBaseBranchLink(validURL)
	assert.NoError(t, err, "Expected no error for a valid URL")
	assert.Equal(t, validURL, config.BaseBranchLink, "BaseBranchLink should be updated to the valid URL")

	// Test with an invalid URL
	invalidURL := "::not a valid url::"
	err = SetBaseBranchLink(invalidURL)
	assert.Error(t, err, "Expected an error for an invalid URL")
	// Check if the link remained unchanged from the previous valid setting
	assert.Equal(t, validURL, config.BaseBranchLink, "BaseBranchLink should not change on invalid URL error")

	// Test with an empty string
	// Reset to default or known state before this
	config.BaseBranchLink = originalLink // Reset to original default
	err = SetBaseBranchLink("")
	assert.NoError(t, err, "Expected no error for an empty string")
	assert.Equal(t, originalLink, config.BaseBranchLink, "BaseBranchLink should not change for an empty string")

	// Clean up global state if necessary for other tests, though each test should manage its own state.
	config = NewConfig()
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Ensure global config is reset at the beginning of this test
	config = NewConfig()

	tempDir := t.TempDir()
	path.SetTestTmpFolder(tempDir) // Configure our path package to use this tempDir for Tmp()
	defer path.ResetTestTmpFolder()  // Clean up path override

	// Modify global config state using exported functions
	testURL := "http://example.com/testbranch"
	err := SetBaseBranchLink(testURL)
	require.NoError(t, err, "SetBaseBranchLink should not fail")

	testDir := "testdir_to_exclude"
	testFile := "testfile_to_exclude.go"
	emptyStr := "" // For non-provided flags

	// Need to pass pointers for Exclude flags
	err = Exclude(&testDir, &testFile, &emptyStr, &emptyStr, &emptyStr, &emptyStr)
	require.NoError(t, err, "Exclude should not fail")

	// Save the modified global config
	// Save uses path.Tmp(Name) internally, which is now redirected
	saveErr := save() // Calling internal save directly to check its error
	require.NoError(t, saveErr, "save() should not fail")

	// Ensure the file was actually created
	savedFilePath := path.Tmp(Name)
	_, statErr := os.Stat(savedFilePath)
	require.NoError(t, statErr, "Config file should exist after save")

	// Create a new, clean Config struct to load into for verification
	loadedCfg := NewConfig()
	// This simulates loading the config from file into a fresh application state
	err = internal.GetAndUnmarshalFile(savedFilePath, loadedCfg)
	require.NoError(t, err, "Loading saved config should not fail")

	// Assertions: Check if the loaded config matches what was saved
	assert.Equal(t, testURL, loadedCfg.BaseBranchLink, "Loaded BaseBranchLink does not match")

	require.NotNil(t, loadedCfg.ExDirs, "Loaded ExDirs map should not be nil")
	assert.True(t, loadedCfg.ExDirs[testDir], "Loaded ExDirs does not contain expected directory")

	require.NotNil(t, loadedCfg.ExFiles, "Loaded ExFiles map should not be nil")
	assert.True(t, loadedCfg.ExFiles[testFile], "Loaded ExFiles does not contain expected file")

	// As an extra check, ensure the global `config` wasn't accidentally used for loading assertions
	// by resetting it before loading into `loadedCfg` and checking `loadedCfg` specifically.
	// The previous `config = NewConfig()` before `internal.GetAndUnmarshalFile` handles this.
}
