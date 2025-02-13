package ops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetRoot(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "testproject")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup after test

	// Create a fake project structure
	rootDir := filepath.Join(tempDir, "myproject")
	subDir := filepath.Join(rootDir, "subfolder")

	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	// Create a marker file (go.mod)
	markerFile := filepath.Join(rootDir, "go.mod")
	if err := os.WriteFile(markerFile, []byte("module testmodule"), 0644); err != nil {
		t.Fatalf("Failed to create marker file: %v", err)
	}

	// Change working directory to the subfolder
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd) // Restore original working directory after test

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Run the getRoot function
	detectedRoot, err := getRoot()
	if err != nil {
		t.Fatalf("getRoot() returned an error: %v", err)
	}

	// Verify that the detected root is correct
	if detectedRoot != rootDir {
		t.Errorf("Expected root directory: %s, got: %s", rootDir, detectedRoot)
	}
}
