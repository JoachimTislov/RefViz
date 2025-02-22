package ops

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/manifoldco/promptui"
)

// Scan scans the content for symbols and references
// If the content is a directory, it scans all files in the directory
// If scanForRefs is true, it scans for references
// If the content is a file, it only scans the file
func Scan(content *string, scanAgain *bool) (*time.Time, error) {
	paths, err := findContent(content)
	if err != nil {
		return nil, fmt.Errorf("error finding content: %s, err: %v", *content, err)
	}
	// Start the timer
	// This is used to calculate the time it takes to scan the content
	startNow := time.Now()
	for _, path := range paths {
		if err := processPath(path, *scanAgain); err != nil {
			return nil, fmt.Errorf("error processing path: %s, err: %v", path, err)
		}
	}
	return &startNow, nil
}

func processPath(path string, scanAgain bool) error {
	e, valid := checkIfValid(path)
	if !valid {
		return fmt.Errorf("error: %s is not a valid entity", path)
	}
	if e.IsDir() {
		return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking through directory: %s, err: %v", path, err)
			}
			if !d.IsDir() && isValid(d.IsDir(), path) {
				return getContent(path, scanAgain)
			}
			return nil
		})
	}
	return getContent(path, scanAgain)
}

// findContent walks for the content root and attempts to find the content
// Returns early if the content is an empty string, equal to scanning everything from project root
func findContent(content *string) ([]string, error) {
	var paths []string
	var err error
	projectRootPath := projectPath()
	if *content == "" {
		paths = append(paths, projectRootPath)
	} else {
		err = filepath.WalkDir(projectRootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking from path: %s, err: %v", path, err)
			}
			_, valid := checkIfValid(path)
			if !valid {
				return nil
			}
			if filepath.Base(path) == *content {
				paths = append(paths, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error finding content: %s, err: %v", *content, err)
		}
	}

	if len(paths) == 0 {
		log.Fatal("No content found")
	}

	paths, err = askUser(paths, []string{})
	if err != nil {
		return nil, fmt.Errorf("error asking user: %v", err)
	}

	return paths, nil
}

func askUser(paths []string, selectedPaths []string) ([]string, error) {
	prompt := pathsPrompt(paths, len(selectedPaths))
	// Run the interactive selection
	index, value, err := prompt.Run()
	if err != nil {
		fmt.Println("Selection failed:", err)
		return nil, err
	}
	if value == scanAll || len(paths) == 1 {
		return paths, nil
	}
	if value == scanSelected {
		return selectedPaths, nil
	}
	if value == exit {
		log.Fatal("Cancelled by user")
	}
	// View createPrompt function to understand the logic
	// Scan all is at index 0 and exit is at index len(paths) - 1
	if index >= 1 && index < len(paths)+1 {
		index -= 1
		paths = append(paths[:index], paths[index+1:]...)
	}
	return askUser(paths, append(selectedPaths, value))
}

func pathsPrompt(paths []string, lenSelectedPaths int) promptui.Select {
	if len(paths) > 1 && lenSelectedPaths == 0 {
		paths = append([]string{scanAll}, paths...)
	}
	if lenSelectedPaths > 0 {
		paths = append([]string{scanSelected}, paths...)
	}
	return selectPrompt("Select content to scan", append(paths, exit))
}
