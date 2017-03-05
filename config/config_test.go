package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig1(t *testing.T) {
	conf, err := FromBytes([]byte(`

target "default" {
  docker {
    build_script = "./dockerbuild.sh"
    dockerfile_build = "Dockerfile-build"
    dockerfile = "Dockerfile"
    image = "localhost:5000/internal-kube1/myimage"
    forwarder = "registry"
  }

  deployment {
    cluster = "data-priv"
    name = "bagpipe-v2-realtime"
    container = "bagpipe"
  }
}


forwarder "registry" {
  cluster = "data-priv"
  service = "registry"
  namespace = "kube-system"
  port = 5000
  local_port = 5000
}

forwarder "dw" {
  pod = "app=goflow-dw"
  port = 7777
}

forwarder "ssh" {
  ssh_host = "my-docker-registry.local"
  port = 5000
}

`))
	require.NoError(t, err)

	require.NotNil(t, conf.Targets)
	require.NotNil(t, conf.Forwarders)

	require.NotNil(t, conf.Targets["default"])
	require.NotNil(t, conf.Targets["default"].Docker)
	assert.Equal(t, conf.Targets["default"].Docker.BuildScript, "./dockerbuild.sh")
	assert.Equal(t, conf.Targets["default"].Docker.DockerfileBuild, "Dockerfile-build")
	assert.Equal(t, conf.Targets["default"].Docker.Dockerfile, "Dockerfile")
	assert.Equal(t, conf.Targets["default"].Docker.Image, "localhost:5000/internal-kube1/myimage")
	assert.Equal(t, conf.Targets["default"].Docker.Forwarder, "registry")

	require.NotNil(t, conf.Targets["default"].Deployment)
	assert.Equal(t, conf.Targets["default"].Deployment.Cluster, "data-priv")
	assert.Equal(t, conf.Targets["default"].Deployment.Name, "bagpipe-v2-realtime")
	assert.Equal(t, conf.Targets["default"].Deployment.Container, "bagpipe")

	require.NotNil(t, conf.Forwarders["registry"])
	assert.Equal(t, conf.Forwarders["registry"].Cluster, "data-priv")
	assert.Equal(t, conf.Forwarders["registry"].Service, "registry")
	assert.Equal(t, conf.Forwarders["registry"].Namespace, "kube-system")
	assert.Equal(t, conf.Forwarders["registry"].Port, 5000)
	assert.Equal(t, conf.Forwarders["registry"].LocalPort, 5000)

	require.NotNil(t, conf.Forwarders["dw"])
	assert.Equal(t, conf.Forwarders["dw"].Pod, "app=goflow-dw")
	assert.Equal(t, conf.Forwarders["dw"].Port, 7777)

	require.NotNil(t, conf.Forwarders["ssh"])
	assert.Equal(t, conf.Forwarders["ssh"].SSHHost, "my-docker-registry.local")
	assert.Equal(t, conf.Forwarders["ssh"].Port, 5000)
}

func TestConfig2(t *testing.T) {
	conf, err := FromBytes([]byte(`

target "default" {
  // Uses the build script in the current directory to build
  docker {
    build_script = "./dockerbuild.sh"
    image = "localhost:5000/internal-kube1/myimage"
    tag = "from-file"
    tag_file = "VERSION.txt"
  }
}

target "without_build_scripts" {
  docker {
    dockerfile_build = "Dockerfile-build"
    dockerfile = "Dockerfile"
    workdir = "./docker"
    image = "localhost:5000/internal-kube1/myimage"
  }
}

target "with_default_values" {
  docker {
    dockerfile_build = "Dockerfile-build" // always used if present
    dockerfile = "Dockerfile" // default value, used directly to build if no "build_script" is specified.
    workdir = "."  // default value
    tag = "git-short-rev" // default value
    image = "localhost:5000/internal-kube1/myimage"
  }
}

target "same_as_previous" {
  docker {
    image = "localhost:5000/internal-kube1/myimage"
  }
}

default_target = "default"  // default value, can be overridden.


`))

	require.NoError(t, err)

	require.NotNil(t, conf.Targets)
	require.Nil(t, conf.Forwarders)

	require.NotNil(t, conf.Targets["default"])
	assert.Len(t, conf.Targets, 4)

	require.NotNil(t, conf.Targets["default"].Docker)
	assert.Equal(t, conf.Targets["default"].Docker.BuildScript, "./dockerbuild.sh")
	assert.Equal(t, conf.Targets["default"].Docker.Tag, "from-file")
	assert.Equal(t, conf.Targets["default"].Docker.TagFile, "VERSION.txt")

	require.NotNil(t, conf.Targets["without_build_scripts"].Docker)
	assert.Equal(t, conf.Targets["without_build_scripts"].Docker.DockerfileBuild, "Dockerfile-build")
	assert.Equal(t, conf.Targets["without_build_scripts"].Docker.Dockerfile, "Dockerfile")
	assert.Equal(t, conf.Targets["without_build_scripts"].Docker.Workdir, "./docker")
	assert.Equal(t, conf.Targets["without_build_scripts"].Docker.Image, "localhost:5000/internal-kube1/myimage")

	require.NotNil(t, conf.Targets["with_default_values"].Docker)
	assert.Equal(t, conf.Targets["with_default_values"].Docker.Tag, "git-short-rev")

	require.NotNil(t, conf.Targets["same_as_previous"].Docker)
	assert.Equal(t, conf.Targets["same_as_previous"].Docker.Image, "localhost:5000/internal-kube1/myimage")

	assert.Equal(t, conf.DefaultTarget, "default")
}

func TestValidate(t *testing.T) {
	conf := &RepoConfig{
		DefaultTarget: "doesn't-exist",
	}
	assert.EqualError(t, conf.Validate(), "2 errors occurred:\n\n* no sections defined, expecting \"target\" and/or \"forwarder\" statements\n* default_target \"doesn't-exist\" specified points to an undefined target")

	conf = &RepoConfig{
		Targets: targetMap("default", &Target{Docker: &Docker{Tag: "not-supported"}}),
	}
	assert.EqualError(t, conf.Validate(), "2 errors occurred:\n\n* target \"not-supported\": tag \"default\" invalid, options are \"from-file\" and \"git-short-rev\"\n* target \"default\": docker \"image\" field required")

	conf = &RepoConfig{
		Targets: targetMap("default", &Target{Docker: &Docker{Image: "img", Forwarder: "not-defined"}}),
	}
	assert.EqualError(t, conf.Validate(), "1 error occurred:\n\n* target \"default\": forwarder \"not-defined\" not defined")

	conf = &RepoConfig{
		Targets: targetMap("default", &Target{Docker: &Docker{Image: "img", Tag: "from-file", TagFile: ""}}),
	}
	assert.EqualError(t, conf.Validate(), "1 error occurred:\n\n* target \"default\": \"tag_file\" missing (as \"tag\" is \"from-file\")")

}

func targetMap(elements ...interface{}) map[string]*Target {
	out := make(map[string]*Target)
	for i := 0; i < len(elements); i += 2 {
		key := elements[i].(string)
		val := elements[i+1].(*Target)
		out[key] = val
	}
	return out
}
