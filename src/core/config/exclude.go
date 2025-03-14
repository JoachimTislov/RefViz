package config

import (
	"fmt"
)

func exclude(m *SbMap, items ...string) error {
	for i := range items {
		_, valid := CheckIfValid(items[i])
		if valid {
			(*m)[items[i]] = true
		} else {
			return fmt.Errorf("invalid exclusion: %s", items[i])
		}
	}
	return config.save()
}
