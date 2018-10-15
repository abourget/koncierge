package build

import (
	"context"
	"fmt"
	"os/exec"
)

func (b *Build) Forward(forwarder string) (kill func(), err error) {
	fwd := b.Config.Forwarders[forwarder]
	if fwd == nil {
		return nil, fmt.Errorf("forwarder %q unknown", forwarder)
	}

	ctx, cancel := context.WithCancel(context.Background())
	kill = cancel

	_ = exec.CommandContext(ctx, "ssh", "-L", "1234:localhost:2354")

	return
}
