package prompt

import (
	"fmt"
	"log"
	"slices"

	"github.com/manifoldco/promptui"
)

const (
	scanAll      = "Scan all content"
	exit         = "Cancel"
	scanSelected = "Scan selected content"
)

func AskUser(paths []string, selectedPaths []string) ([]string, error) {
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
		paths = slices.Delete(paths, index, index+1)
	}
	return AskUser(paths, append(selectedPaths, value))
}

func pathsPrompt(paths []string, lenSelectedPaths int) promptui.Select {
	if len(paths) > 1 && lenSelectedPaths == 0 {
		paths = append([]string{scanAll}, paths...)
	}
	if lenSelectedPaths > 0 {
		paths = append([]string{scanSelected}, paths...)
	}
	return SelectPrompt("Select content to scan", append(paths, exit))
}
