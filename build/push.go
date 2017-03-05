package build

import (
	"fmt"
	"os"
	"os/exec"
)

func (b *Build) Push(target string) error {
	t := b.Config.Targets[target]

	tag, err := b.getTag(t)
	if err != nil {
		return fmt.Errorf("could not get tag: %s", err)
	}

	// Check authentication before we go, or re-authenticate before we go, or CHECK that!

	imageTag := fmt.Sprintf("%s:%s", t.Image, tag)
	cmd := exec.Command("docker", "push", imageTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("koncierge: pushing docker image %q\n", imageTag)

	return cmd.Run()
}
