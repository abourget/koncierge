package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/homedir"
)

func RunKubectl(cluster string, namespace string, args []string) error {

	// TODO find  `~/.kube/[cluster]/kubeconfig|config`,
	// craft the `kubectl` command, adding --namespace` if namespace != ""
	// append the other params,
	// set the right `KUBECONFIG` according to cluster configuration..
	kubeconfPath := filepath.Join(homedir.Get(), ".kube", cluster, "config")
	if _, err := os.Stat(kubeconfPath); err != nil {
		return err
	}

	cmd := exec.Command("kubectl", "--kubeconfig", kubeconfPath)
	if namespace != "" {
		cmd.Args = append(cmd.Args, "--namespace", namespace)
	}
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("koncierge: running %q\n", cmdtostring(cmd))

	return cmd.Run()
}

func cmdtostring(c *exec.Cmd) string {
	out := c.Path
	args := strings.Join(c.Args[1:], " ")
	if args != "" {
		out = fmt.Sprintf("%s %s", out, args)
	}
	return out
}
