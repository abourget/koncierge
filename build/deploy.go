package build

import "fmt"

func (b *Build) Deploy(target string) error {
	t := b.Config.Targets[target]

	tag, err := b.getTag(t)
	if err != nil {
		return fmt.Errorf("could not get tag: %s", err)
	}

	// TODO: find the `kubeconfig` file or `config` file in `~/.kube/[name-of-cluster]/`..
	// use that as `KUBECONFIG` when calling `kubectl` (which must be in $PATH),
	//  do a `set image` with all the details
	// we'll want to re-use this configuration setup when doing port-forwards,
	// as we'll then again, call `kubectl` to get some pod values,
	// port-forward calls, etc..
	imageTag := fmt.Sprintf("%s:%s", t.Image, tag)
	return RunKubectl(t.Deployment.Cluster, t.Deployment.Namespace, []string{
		"set", "image", t.Deployment.Name, fmt.Sprintf("%s=%s", t.Deployment.Container, imageTag),
	})
}
