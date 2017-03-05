package build

import (
	"fmt"
	"log"
)

func (b *Build) Push(target string) error {
	t := b.Config.Targets[target]

	tag, err := b.getTag(t)
	if err != nil {
		return fmt.Errorf("could not get tag: %s", err)
	}

	log.Println("Pushing to Docker, implement me", tag, t.DockerfileWithDefault())
	return nil
}
