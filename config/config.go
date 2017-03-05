package config

import (
	"fmt"
	"io/ioutil"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
)

type RepoConfig struct {
	Targets       map[string]*Target    `hcl:"target"`
	Forwarders    map[string]*Forwarder `hcl:"forwarder"`
	DefaultTarget string                `hcl:"default_target"`

	// FilePath is the full path where the configuration was loaded from.
	FilePath string `hcl:"-"`
}

type Deployment struct {
	Cluster   string `hcl:"cluster"`
	Namespace string `hcl:"namespace"`
	Name      string `hcl:"name"`
	Container string `hcl:"container"`
}

type Forwarder struct {
	Cluster   string `hcl:"cluster"`
	Namespace string `hcl:"namespace"`
	Service   string `hcl:"service"`
	Pod       string `hcl:"pod"`

	SSHHost string `hcl:"ssh_host"`

	Port      int `hcl:"port"`
	LocalPort int `hcl:"local_port"`
}

func FromFile(filename string) (*RepoConfig, error) {
	rawConf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	conf, err := FromBytes(rawConf)
	if err != nil {
		return nil, err
	}

	conf.FilePath = filename

	return conf, nil
}

// FromBytes reads the HCL configuration from a byte-slice.
func FromBytes(config []byte) (*RepoConfig, error) {
	var c RepoConfig
	err := hcl.Unmarshal(config, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *RepoConfig) Validate() (out error) {
	if c.Targets == nil && c.Forwarders == nil {
		out = multierror.Append(out, fmt.Errorf(`no sections defined, expecting "target" and/or "forwarder" statements`))
	}

	if c.DefaultTarget != "" && (c.Targets == nil || c.Targets[c.DefaultTarget] == nil) {
		out = multierror.Append(out, fmt.Errorf(`default_target %q specified points to an undefined target`, c.DefaultTarget))
	}
	if c.Targets != nil {
		for targetName, target := range c.Targets {
			tag := target.Tag
			switch tag {
			case "git-short-rev", "from-file", "":
			default:
				out = multierror.Append(out, fmt.Errorf(`target %q: tag %q invalid, options are "from-file" and "git-short-rev"`, tag, targetName))
			}

			if tag == "from-file" && target.TagFile == "" {
				out = multierror.Append(out, fmt.Errorf(`target %q: "tag_file" missing (as "tag" is "from-file")`, targetName))
			}

			forwarder := target.Forwarder
			if forwarder != "" && (c.Forwarders == nil || c.Forwarders[forwarder] == nil) {
				out = multierror.Append(out, fmt.Errorf(`target %q: forwarder %q not defined`, targetName, forwarder))
			}

			if target.Image == "" {
				out = multierror.Append(out, fmt.Errorf(`target %q: docker "image" field required`, targetName))
			}

			if target.Deployment != nil {
				if target.Deployment.Cluster == "" {
					out = multierror.Append(out, fmt.Errorf(`target %q: deployment's "cluster" statement missing`, targetName))
				}
				if target.Deployment.Name == "" {
					out = multierror.Append(out, fmt.Errorf(`target %q: deployment's "name" statement missing`, targetName))
				}
				if target.Deployment.Container == "" {
					out = multierror.Append(out, fmt.Errorf(`target %q: deployment's "container" statement missing`, targetName))
				}
			}
		}
	}
	if c.Forwarders != nil {
		for fwderName, fwder := range c.Forwarders {
			if fwder.SSHHost != "" && fwder.Cluster != "" {
				out = multierror.Append(out, fmt.Errorf(`forwarder %q: mutually exclusive "ssh_host" and "cluster" specified`, fwderName))
			}

			if fwder.SSHHost == "" && fwder.Cluster == "" {
				out = multierror.Append(out, fmt.Errorf(`forwarder %q: "ssh_host" or "cluster" required, what will this forwarder do otherwise ?`, fwderName))
			}

			if fwder.Cluster == "" {
				if fwder.Pod != "" {
					out = multierror.Append(out, fmt.Errorf(`forwarder %q: "cluster" required if "pod" present`, fwderName))
				}
				if fwder.Service != "" {
					out = multierror.Append(out, fmt.Errorf(`forwarder %q: "cluster" required if "service" present`, fwderName))
				}
				if fwder.Namespace != "" {
					out = multierror.Append(out, fmt.Errorf(`forwarder %q: "cluster" required if "namespace" present`, fwderName))
				}
			}

			if fwder.Port == 0 {
				out = multierror.Append(out, fmt.Errorf(`forwarder %q: "port" is required`, fwderName))
			}
		}
	}

	return
}
