package content

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JoachimTislov/RefViz/core/config"
	"github.com/JoachimTislov/RefViz/internal/routines"
)

// Scan scans the content for symbols and references
// If the content is a directory, it scans all files in the directory
// If scanForRefs is true, it scans for references
// If the content is a file, it only scans the file
func Scan(content *string, scanAgain, ask *bool) error {
	paths, err := Find(content, ask)
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
	e, valid := config.CheckIfValid(path)
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
		jobs = append(jobs, Get(path, scanAgain, everythingIsUpToDate))
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
		if !d.IsDir() && config.IsValid(d.IsDir(), path) {
			*paths = append(*paths, path)
		}
		return nil
	})
}
