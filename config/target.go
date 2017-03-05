package config

type Target struct {
	Workdir string `hcl:"workdir"`

	BuildScript string `hcl:"build_script"`

	// The Dockerfile is relative to the Workdir
	Dockerfile string `hcl:"dockerfile"`

	Auth    string `hcl:"auth"` // docker push authentication..
	Image   string `hcl:"image"`
	Tag     string `hcl:"tag"` // 'from-file', 'git-short-rev'
	TagFile string `hcl:"tag_file"`

	Forwarder string `hcl:"forwarder"`

	Deployment *Deployment `hcl:"deployment"`
}

func (t *Target) DockerfileWithDefault() string {
	if t.Dockerfile != "" {
		return t.Dockerfile
	}
	return "Dockerfile"
}
