package helpers

import (
	"fmt"
	"os"
	"testing"

	"github.com/JoachimTislov/RefViz/core/cache"
	"github.com/JoachimTislov/RefViz/core/load"
	"github.com/JoachimTislov/RefViz/internal/path"
	"github.com/JoachimTislov/RefViz/internal/utils/rand"
)

// ExecuteTestSequence loads the test cache, runs the tests, and cleans up the test cache
// Should be called in the TestMain function of the test file to ensure clean test cache
func ExecuteTestSequence(m *testing.M) {
	c := loadTestCache()
	m.Run()
	cleanupTestCache(c)
}

func loadTestCache() string {
	t := fmt.Sprintf("test-cache%s.json", rand.String())
	if err := load.File(path.Tmp(t), cache.Get()); err != nil {
		fmt.Printf("error loading cache: %v", err)
	}
	return t
}

func cleanupTestCache(testCache string) {
	if err := os.Remove(path.Tmp(testCache)); err != nil {
		fmt.Printf("error removing cache: %v", err)
	}
}
