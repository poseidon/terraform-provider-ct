# ct_config Data Source

Validate a [Butane config](https://coreos.github.io/butane/specs/) and transpile it to an [Ignition config](https://coreos.github.io/ignition/) for machine consumption.

## Usage

```hcl
data "ct_config" "worker" {
  content      = file("worker.yaml")
  strict       = true
  pretty_print = false

  snippets = [
    file("units.yaml"),
    file("storage.yaml"),
  ]
}

resource "aws_instance" "worker" {
  user_data = data.ct_config.worker.rendered
}
```

See the [Flatcar Linux](../../examples/flatcar-linux.tf), [Fedora CoreOS](../../examples/fedora-coreos.tf), [Fedora IoT](../../examples/fedora-iot.tf), [OpenShift](../../examples/openshift.tf), or [RHEL for Edge](../../examples/rhel-edge.tf) examples.

## Argument Reference

* `content` - contents of a Butane Config that should be validated and transpiled to Ignition (supports `fcos`, `flatcar`, `fiot`, `r4e`, and `openshift` variants).
* `strict` - strictly treat validation warnings as errors (default: false).
* `pretty_print` - indent transpiled Ignition for visual prettiness (default: false).
* `use_mapped_version` - if true, output the Ignition version corresponding to the given Butane version, matching Butane CLI behavior (default: false).
* `files_dir` - allow embedding local files relative to this directory.
* `snippets` - list of Butane snippets to merge into the content. Snippets must have the same `variant` as the main content.

## Argument Attributes

* `rendered` - rendered Ignition configuration (or MachineConfig YAML for the `openshift` variant).

## Versions

Butane configs and snippets are automatically upgraded to the latest stable Ignition specification supported by the provider.
