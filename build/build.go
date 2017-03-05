package build

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/abourget/koncierge/config"
)

type Build struct {
	Config    *config.RepoConfig
	CachedTag string
}

func New(conf *config.RepoConfig) *Build {
	return &Build{
		Config: conf,
	}
}

func (b *Build) Build(target string) error {
	t := b.Config.Targets[target]

	tag, err := b.getTag(t)
	if err != nil {
		return fmt.Errorf("could not get tag: %s", err)
	}

	imageTag := fmt.Sprintf("%s:%s", t.Image, tag)
	workdir := formatWorkdir(b.Config, t)
	env := os.Environ()
	env = append(env, fmt.Sprintf("KONCIERGE_IMAGE=%s", t.Image))
	env = append(env, fmt.Sprintf("KONCIERGE_TAG=%s", tag))
	env = append(env, fmt.Sprintf("KONCIERGE_IMAGE_TAG=%s", imageTag))

	if t.BuildScript != "" {
		cmd := exec.Command(t.BuildScript)
		cmd.Dir = workdir
		cmd.Env = env
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("koncierge: build script %q failed: %s\n", t.BuildScript, err)
		}
	}

	cmd := exec.Command("docker", "build", "-t", imageTag, "-f", t.DockerfileWithDefault(), ".")
	cmd.Dir = workdir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("koncierge: docker build command failed:", err)
	}

	return nil
}

func (b *Build) TargetWithDefault(cliTarget string) (string, error) {
	target := b.defaultTarget(cliTarget)

	if b.Config.Targets[target] == nil {
		return target, fmt.Errorf("target %q is not defined", target)
	}
	return target, nil
}
func (b *Build) defaultTarget(cliTarget string) string {
	if cliTarget != "" {
		return cliTarget
	}
	if b.Config.DefaultTarget != "" {
		return b.Config.DefaultTarget
	}
	return "default"
}