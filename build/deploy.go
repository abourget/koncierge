package build

import (
	"fmt"
	"log"
)

func (b *Build) Deploy(target string) error {
	t := b.Config.Targets[target]

	tag, err := b.getTag(t)
	if err != nil {
		return fmt.Errorf("could not get tag: %s", err)
	}

	log.Println("Deploying to Kubernetes, implement me", tag, t.DockerfileWithDefault())
	return nil
}
