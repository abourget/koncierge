package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// WalkConfig walks up the parent directories (starting from the
// current workdir), until it finds a `.koncierge` or a
// `Konciergefile` file, it then attempts to load that file.
func WalkConfig() (c *RepoConfig, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	var konciergeFilePath string
	var previousPath string
	for {
		konciergeFilePath = filepath.Join(cwd, "Konciergefile")
		if _, err := os.Stat(konciergeFilePath); err == nil {
			break
		}

		nextPath := filepath.Dir(cwd)
		if nextPath == previousPath {
			return nil, fmt.Errorf("not a koncierge project (or any of the parent directories), searching for Konciergefile files")
		}
		previousPath = nextPath
	}

	return FromFile(konciergeFilePath)
}
