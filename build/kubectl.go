package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunKubectl(ctx context.Context, cluster string, namespace string, args []string) error {
	envVar := fmt.Sprintf("KONCIERGE_%s_KUBECONFIG", strings.Replace(strings.ToUpper(cluster), "-", "_", -1))
	kubeconfPath := os.Getenv(envVar)
	if kubeconfPath == "" {
		return fmt.Errorf("couldn't find environment variable %q pointing to a KUBECONFIG file", envVar)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfPath)
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
