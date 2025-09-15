package load

import (
	"fmt"

	"github.com/JoachimTislov/refviz/core/cache"
	"github.com/JoachimTislov/refviz/core/config"
	"github.com/JoachimTislov/refviz/internal/path"
)

func Defs() error {
	_ = path.Project()

	if err := loadConfig(config.Name); err != nil {
		return fmt.Errorf("error loading configurations: %v", err)
	}
	if err := loadCache(cache.Name); err != nil {
		return fmt.Errorf("error loading cache: %v", err)
	}
	return nil
}
