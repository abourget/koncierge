package build

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/abourget/koncierge/config"
)

func (b *Build) getTag(t *config.Target) (string, error) {
	if b.CachedTag != "" {
		return b.CachedTag, nil
	}
	return b.getUncachedTag(t)
}

func (b *Build) getUncachedTag(t *config.Target) (string, error) {
	workdir := formatWorkdir(b.Config, t)

	switch t.Tag {
	case "", "git-short-rev":
		cmd := exec.Command("git", "describe", "--long", "--always", "--dirty")
		cmd.Dir = workdir
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(out)), nil
	case "tag-file":
		cnt, err := ioutil.ReadFile(filepath.Join(t.Workdir, t.TagFile))
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(cnt)), nil
	}
	return "", fmt.Errorf("unsupported tag method: %q", t.Tag)
}
