package config

func AddExDirs(dirs ...string) error {
	return exclude(&config.ExDirs, dirs...)
}

func AddExFiles(files ...string) error {
	return exclude(&config.ExDirs, files...)
}
