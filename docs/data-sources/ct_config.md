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

See the [Flatcar Container Linux](../../examples/flatcar-linux.tf) or [Fedora CoreOS](../../examples/fedora-coreos.tf) examples.

## Argument Reference

* `content` - contents of a Butane Config that should be validated and transpiled to Ignition.
* `strict` - strictly treat validation warnings as errors (default: false).
* `pretty_print` - indent transpiled Ignition for visual prettiness (default: false)
* `snippets` - list of Butane snippets to merge into the content. Content and snippet configs must have the same `version` and `variant`.

## Argument Attributes

* `rendered` - transpiled Ignition configuration

