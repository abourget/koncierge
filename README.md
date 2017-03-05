Koncierge, a Kubernetes + Containers developer concierge
========================================================

Koncierge helps you with the task of building, pushing and deploying
containers backed by a Kubernetes cluster.

It also helps distribute to developers credentials necessary to access
the different available clusters, in a flow that is both easy and secure.

There are two major feature sets of Koncierge:

1. Read and understand a `Konciergefile` file in a Git repository, and
   ease operations on such a repo / project.
2. Manage cluster configurations in `~/.kube` directory, and ease
   operations on the cluster configuration.

The codebase tries to mimic some `kubernetes` Go idioms, like the
directory layout

## Example Koncierge repositories

### First example

With this `Konciergefile` file:

```hcl

target "default" {
  build_script = "./dockerbuild.sh"
  image = "123123123123213.dkr.ecr.amazonaws.com/data-priv/myproject"
  auth = "aws-ecr"
  forwarder = "ssh"

  deployment {
    cluster = "data-priv"
    name = "bagpipe-v2-realtime"
    namespace = "default"
    container = "bagpipe"
  }
}

forwarder "ssh" {
  ssh_host = "my-docker-registry.local"
  port = 5000
}

```

someone can run, any of these commands:

```bash
koncierge build --push --deploy
koncierge build
koncierge -t default build
koncierge push
koncierge deploy
koncierge fwd ssh -- docker push localhost:5000/other/image
koncierge logs -f
```

If you specify a `build_script`, it is executed

### Second example

With such a `Konciergefile` file:


```hcl

target "default" {
  build_script = "./dockerbuild.sh"
  image = "localhost:5000/internal-kube1/myimage"
  forwarder = "registry"

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
  cluster = "data-priv"
  pod = "app=goflow-dw"
  port = 7777
}
```

One can run those commands fruitfully:

```
koncierge fwd dw
koncierge build --push --deploy
koncierge push --deploy
koncierge deploy
koncierge print fwd dw pod  # prints `goflow-dw-12345123412-eifje`
koncierge print fwd dw namespace  # prints `default`
koncierge -t second_target print image
koncierge print deployment name
```

### Third example - build steps

```hcl
target "default" {
  // Uses the build script in the current directory to build
  build_script = "./dockerbuild.sh"
  image = "localhost:5000/internal-kube1/myimage"
  tag = "from-file"
  tag_file = "VERSION.txt"
}

target "without_build_scripts" {
  dockerfile = "Dockerfile"
  workdir = "./docker"
  image = "localhost:5000/internal-kube1/myimage"
}

target "with_default_values" {
  dockerfile = "Dockerfile" // default value, used directly to build if no `build_script` is specified.
  workdir = "."  // default value
  tag = "git-short-rev" // default value
  image = "localhost:5000/internal-kube1/myimage"
}

target "same_as_previous" {
  image = "localhost:5000/internal-kube1/myimage"
}

default_target = "default"  // default value, can be overridden.

```

When a `build_script` is defined, the following environment variables
are injected in the child process:

* `KONCIERGE_IMAGE`, an image name, based on the current config.
* `KONCIERGE_TAG`, a tag value, based on the current `tag` algorithm.

### Tag algorithms

The following `tag` algorithms are available to determine the tag of
the built Docker images.

* `git-short-rev`, runs a `git describe --long --always --dirty` on
  the repo and provides that value as tag.
* `from-file`, reads the `.tag` file in the current directory.
  Filename can be overridden with the `tag_file` instruction.  If a
  `build_script` is specified alongside this mode, the `KONCIERGE_TAG`
  will be the value previously present in that file, but the file will
  be re-read after the build script has executed (in case the script
  updated it).

## Discovery of the `Konciergefile` file

Directories are searched starting from the current directory, going up
until hitting the $HOME path.  The first (deepest) is taken and
used. The current working directory is also changed (for `koncierge
build`) to the directory where the `Konciergefile` file is defined.


## Auto-detection of environment

Koncierge detects if you're on Windows, Mac or Linux, and will tweak
its configuration and automatically load `docker-machine env`, so you
can always call it without first starting any environment-wrapping
`eval $(docker-machine env)` calls.

You can override the autodetection with the `KONCIERGE_DETECT_DOCKER`
set to `0`.


# Koncierge as an authentication wrapper

Koncierge knows about your Docker registries and Kubernetes
clusters. It can help you authenticate with them, or request access to
some new Kubernetes clusters that exist within your organisation.

    koncierge access request https://clusters.truekey.engineering/data-priv
    (report: {"name": "auth_method": "csr", "groups": ["admin", "user", "namespace-sentry", "namespace-kube-system"})
    Username: abourget
    Group:
    (1) admin
    (2) user
    (3) namespace-sentry
    (4) namespace-kube-system
    Request group: 1
    --
    Send to bob@intel.com / send to @hubert on Slack.
    ------ BASE64 BLOB ------
    {"username": "abourget",
     "grantee_gpg_key": "...........",
     "group": "....."

       or

     "csr": "..." (contains CN=abourget, O=groupname, O=groupname2}
    ------ END BASE64 BLOB ------
    (wrote ~/.kube/data-priv/config the CSR and private key)

Auth methods:

 * `csr`, transmits the CSR, receives a certificate
 * `keys`, transmits a temporary gpg key, receives a private/public key pair, encrypted with the temp gpg key


```bash

koncierge access request
koncierge access analyze requestfile.json (from the admin)
koncierge access sign requestfile.json --options (from the admin)
```

The process for an admin to `grant` is specific to your
infrastructure, but `koncierge` will help you manage your different
credentials and load them only when requested by some `Konciergefile`
files.

```bash

$ koncierge env data-priv
PS1=[kube:data-priv] $PS1
```

You can add that to your `.bashrc` or `.zshrc`
