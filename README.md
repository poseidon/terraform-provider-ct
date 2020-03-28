# terraform-provider-ct

`terraform-provider-ct` allows Terraform to validate a [Container Linux Config](https://github.com/coreos/container-linux-config-transpiler/blob/master/doc/configuration.md) or [Fedora CoreOS Config](https://github.com/coreos/fcct/blob/master/docs/configuration-v1_0.md) and render it as [Ignition](https://github.com/coreos/ignition) for machine consumption.

Define a Container Linux Config (CLC) or Fedora CoreOS Config (FCC).

```yaml
# Container Linux Config
---
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

```yaml
# Fedora CoreOS Config
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

Render the config with Terraform for machine consumption.

```hcl
# Define a data source
data "ct_config" "worker" {
  content      = file("worker.yaml")
  pretty_print = false
  strict       = true
  snippets     = []
}

# Usage: Render the config as Ignition
resource "aws_instance" "worker" {
  user_data = data.ct_config.worker.rendered
}
```

See the [Container Linux](examples/container-linux.tf) or [Fedora CoreOS](examples/fedora-coreos.tf) examples.

## Requirements

* Terraform v0.12+ [installed](https://www.terraform.io/downloads.html)

### Ignition Outputs

Container Linux Configs are coupled with the render tool. For example, all CLCs are rendered in Ignition v2.2.0 format. A future `terraform-provider-ct` release would be needed to bump that version.

Fedora CoreOS Config's contain a `version` that is associated with an Ignition format verison. For example, FCC's with `version: 1.0.0` produce Ignition `3.0.0`. A future `terraform-provider-ct` release would be needed to add support for newer versions, but FCCs could continue specifying `1.0.0` indefintely.

| terraform-provider-ct | Ignition (for CLCs) | Ignition (for FCC) |
|-----------------------|---------------------|--------------------|
| 0.5.x                 | Renders 2.2.0       | FCC 1.0.0 -> Ignition 3.0.0 |
| 0.4.x                 | Renders 2.2.0       | FCC 1.0.0 -> Ignition 3.0.0 |
| 0.3.x                 | Renders 2.2.0       | NA                 |
| 0.2.x                 | Renders 2.0.0       | NA                 |

## Install

Add the `terraform-provider-ct` plugin binary for your system to the Terraform 3rd-party [plugin directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) `~/.terraform.d/plugins`.

```sh
VERSION=v0.5.0
wget https://github.com/poseidon/terraform-provider-ct/releases/download/$VERSION/terraform-provider-ct-$VERSION-linux-amd64.tar.gz
tar xzf terraform-provider-ct-$VERSION-linux-amd64.tar.gz
mv terraform-provider-ct-$VERSION-linux-amd64/terraform-provider-ct ~/.terraform.d/plugins/terraform-provider-ct_$VERSION
```

Terraform plugin binary names are versioned to allow for migrations of managed infrastructure.

```
$ tree ~/.terraform.d/
/home/user/.terraform.d/
└── plugins
    ├── terraform-provider-ct_v0.3.0
    ├── terraform-provider-ct_v0.3.1
    ├── terraform-provider-ct_v0.3.2
    └── terraform-provider-ct_v0.5.0
```

## Usage

Configure the ct provider in a `providers.tf` file.

```hcl
provider "ct" {
  version = "0.5.0"
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

Declare a `ct_config` resource in Terraform. Set the `content` to the contents of a Container Linux Config (CLC) or Fedora CoreOS Config (FCC) that should be validated and rendered as Ignition.

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

Use the `snippets` field to append a list of Container Linux Config (CLC) or Fedora CoreOS Config (FCC) snippets. Use `platform` if [platform-specific](https://github.com/coreos/container-linux-config-transpiler/blob/master/config/platform/platform.go) susbstitution is desired (Container Linux only).

## Development

### Binary

To develop the provider plugin locally, build an executable with Go v1.11+.

```
make
```

### Vendor

Add or update dependencies in `go.mod` and vendor.

```
make update
make vendor
```
