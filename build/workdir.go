package build

import (
	"path/filepath"

	"github.com/abourget/koncierge/config"
)

func formatWorkdir(c *config.RepoConfig, t *config.Target) string {
	return filepath.Join(filepath.Dir(c.FilePath), t.Workdir)
}
