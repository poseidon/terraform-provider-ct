# ct_config Data Source

Validate a [Container Linux Config](https://github.com/coreos/container-linux-config-transpiler/blob/master/doc/configuration.md) or a [Butane config](https://coreos.github.io/butane/specs/) and transpile it to an [Ignition config](https://coreos.github.io/ignition/) for machine consumption.

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

See the [Container Linux](examples/container-linux.tf) or [Fedora CoreOS](examples/fedora-coreos.tf) examples.

## Argument Reference

* `content` - contents of a Container Linux Config (CLC) or Fedora CoreOS Config (FCC) that should be validated and transpiled to Ignition.
* `strict` - strictly treat validation warnings as errors (default: false).
* `pretty_print` - indent transpiled Ignition for visual prettiness (default: false)
* `snippets` - list of Container Linux Config (CLC) or Fedora CoreOS Config (FCC) snippets to merge into the content. For FCCs, content and snippet configs must have the same `version`.
* `platform` - (Container Linux only) - enable [platform-specific](https://github.com/coreos/container-linux-config-transpiler/blob/master/config/platform/platform.go) susbstitutions

## Argument Attributes

* `rendered` - transpiled Ignition configuration

