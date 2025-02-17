package ops

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

// Scan scans the content for symbols and references
// If the content is a directory, it scans all files in the directory
// If scanForRefs is true, it scans for references
// If the content is a file, it only scans the file
func Scan(scanAgain *bool, findRefs *bool, content *string, ask bool) error {
	paths, err := findContent(content)
	if err != nil {
		return fmt.Errorf("error finding content: %s, err: %v", *content, err)
	}
	if ask {
		paths, err = askUser(paths, []string{})
		if err != nil {
			return fmt.Errorf("error asking user: %v", err)
		}
	}
	for _, p := range paths {
		e, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("error: %s is not a valid entity, err: %v", p, err)
		}
		if e.IsDir() {
			if err := filepath.WalkDir(p, walk(scanAgain, findRefs)); err != nil {
				return fmt.Errorf("error walking through directory: %s, err: %v", p, err)
			}
		}
		if err := getContent(p, findRefs, scanAgain); err != nil {
			return fmt.Errorf("error getting content: %s, err: %v", p, err)
		}
	}
	return nil
}

// walk walks through the directory and scans the files
func walk(scanAgain *bool, findRefs *bool) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", path, err)
		}
		if !d.IsDir() && isValid(d.IsDir(), path) {
			if err := getContent(path, findRefs, scanAgain); err != nil {
				return fmt.Errorf("error getting content: %s, err: %v", path, err)
			}
		}
		return nil
	}
}

// findContent walks for the content root and attempts to find the content
// Returns early if the content is an empty string, equal to scanning everything from project root
func findContent(content *string) ([]string, error) {
	projectRootPath := projectPath()
	if *content == "" {
		return []string{projectRootPath}, nil
	}
	var paths []string
	err := filepath.WalkDir(projectRootPath, func(path string, d fs.DirEntry, err error) error {
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
	return paths, nil
}

func askUser(paths []string, selectedPaths []string) ([]string, error) {
	if len(paths) == 0 {
		fmt.Println("No content found")
		return nil, fmt.Errorf("no content found")
	}
	prompt := createPrompt(paths)
	// Run the interactive selection
	index, value, err := prompt.Run()
	if err != nil {
		fmt.Println("Selection failed:", err)
		return nil, err
	}
	if value == scanAll || len(paths) == 1 {
		return paths, nil
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

func createPrompt(paths []string) promptui.Select {
	if len(paths) > 1 {
		paths = append([]string{scanAll}, paths...)
	}
	paths = append(paths, exit)
	return promptui.Select{
		Label: "Select content to scan",
		Items: paths,
	}
}
