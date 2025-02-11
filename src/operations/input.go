package op

import (
	"fmt"
	"os"
)

func HandleContentInput(isDir *bool, content *string) error {
	entity, err := os.Stat(*content)
	*isDir = entity.IsDir()
	if err != nil {
		return fmt.Errorf("error analyzing content: %s, err: %v", *content, err)
	}
	if err := isValid(*isDir, content); err != nil {
		return fmt.Errorf("error: %s is not a valid entity, err: %v", *content, err)
	}
	return nil
}
