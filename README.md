# terraform-provider-ct [![Build Status](https://github.com/poseidon/terraform-provider-ct/workflows/test/badge.svg)](https://github.com/poseidon/terraform-provider-ct/actions?query=workflow%3Atest+branch%3Amaster)

`terraform-provider-ct` allows Terraform to validate a [Container Linux Config](https://github.com/coreos/container-linux-config-transpiler/blob/master/doc/configuration.md) or [Fedora CoreOS Config](https://github.com/coreos/fcct/blob/master/docs/configuration-v1_1.md) and transpile it as [Ignition](https://github.com/coreos/ignition) for machine consumption.

## Usage

Configure the config transpiler provider (e.g. `providers.tf`).

```hcl
provider "ct" {}

terraform {
  required_providers {
    ct = {
      source  = "poseidon/ct"
      version = "0.8.0"
    }
  }
}
```

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
version: 1.3.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

Define a `ct_config` data source and render for machine consumption.

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

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

## Requirements

* Terraform v0.13+ [installed](https://www.terraform.io/downloads.html)

## Versions

Fedora CoreOS Config's contain a `version` that is associated with an Ignition format verison. For example, a FCC with `version: 1.0.0` would produce Ignition `3.0.0`, across future releases.

Container Linux Configs render a fixed Ignition version, depending on the `terraform-provider-ct` release, so updating alters the rendered Ignition version.

| terraform-provider-ct | CLC to Ignition     | FCC to Ignition    |
|-----------------------|---------------------|--------------------|
| 0.8.x                 | Renders 2.3.0       | FCC (1.0, 1.1, 1.2, 1.3) -> Ignition (3.0, 3.1, 3.2, 3.2)
| 0.7.x                 | Renders 2.3.0       | FCC (1.0, 1.1, 1.2) -> Ignition (3.0, 3.1, 3.2) |
| 0.6.x                 | Renders 2.3.0       | FCC 1.0.0 -> Ignition 3.0.0, FCC 1.1.0 -> Ignition v3.1.0 |
| 0.5.x                 | Renders 2.2.0       | FCC 1.0.0 -> Ignition 3.0.0 |
| 0.4.x                 | Renders 2.2.0       | FCC 1.0.0 -> Ignition 3.0.0 |
| 0.3.x                 | Renders 2.2.0       | NA                 |
| 0.2.x                 | Renders 2.0.0       | NA                 |

Notes:

* Fedora CoreOS Config `snippets` must match the version set in the content. Version skew among snippets is **not** supported.

## Development

### Binary

To develop the provider plugin locally, build an executable with Go v1.13+.

```
make
```

### Vendor

Add or update dependencies in `go.mod` and vendor.

```
make update
make vendor
```

