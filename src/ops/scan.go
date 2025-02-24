package ops

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JoachimTislov/RefViz/routines"
	"github.com/manifoldco/promptui"
)

// Scan scans the content for symbols and references
// If the content is a directory, it scans all files in the directory
// If scanForRefs is true, it scans for references
// If the content is a file, it only scans the file
func Scan(content *string, scanAgain *bool) error {
	paths, err := findContent(content)
	if err != nil {
		return fmt.Errorf("error finding content: %s, err: %v", *content, err)
	}
	// Start the timer
	// This is used to calculate the time it takes to scan the content
	startNow := time.Now()

	everythingIsUpToDate := true
	for _, path := range paths {
		if err := processPath(path, *scanAgain, &everythingIsUpToDate); err != nil {
			return fmt.Errorf("error processing path: %v", err)
		}
	}

	if everythingIsUpToDate {
		log.Println("Everything is up to date\n Use the -a flag to scan again")
	} else {
		log.Printf("Scan time: %v\n", time.Since(startNow))
	}

	return nil
}

func processPath(path string, scanAgain bool, everythingIsUpToDate *bool) error {
	e, valid := checkIfValid(path)
	if !valid {
		return fmt.Errorf("error: %s is not a valid entity", path)
	}
	paths := []string{path}
	// If the path is a directory, get all the files in the directory
	if e.IsDir() {
		if err := getContentInDir(path, &paths); err != nil {
			return fmt.Errorf("error getting content in directory: %s, err: %v", path, err)
		}
	}

	var jobs []func() error
	for _, path := range paths {
		jobs = append(jobs, getContent(path, scanAgain, everythingIsUpToDate))
	}

	// Start the work for 4 workers
	routines.StartWork(3, jobs)

	return nil
}

func getContentInDir(path string, paths *[]string) error {
	return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through directory: %s, err: %v", path, err)
		}
		if !d.IsDir() && isValid(d.IsDir(), path) {
			*paths = append(*paths, path)
		}
		return nil
	})
}

// findContent walks for the content root and attempts to find the content
// Returns early if the content is an empty string, equal to scanning everything from project root
func findContent(content *string) ([]string, error) {
	var paths []string
	var err error
	projectRootPath := projectPath()
	switch {
	case *content == "":
		paths = append(paths, projectRootPath)
	case filepath.IsAbs(*content):
		paths = append(paths, *content)
	default:
		err = filepath.WalkDir(projectRootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking from path: %s, err: %v", path, err)
			}
			if filepath.Base(path) == *content && isValid(d.IsDir(), path) {
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
	if len(paths) == 1 {
		return paths, nil
	}

	/*paths, err = askUser(paths, []string{})
	if err != nil {
		return nil, fmt.Errorf("error asking user: %v", err)
	}*/

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
	if value == exit {
		log.Fatal("Cancelled by user")
	}
	if value == scanSelected {
		return selectedPaths, nil
	}
	if value == scanAll {
		return paths, nil
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
