# terraform-provider-ct [![Build Status](https://github.com/poseidon/terraform-provider-ct/workflows/test/badge.svg)](https://github.com/poseidon/terraform-provider-ct/actions?query=workflow%3Atest+branch%3Amaster)

`terraform-provider-ct` allows Terraform to validate a [Butane config](https://coreos.github.io/butane/specs/) and transpile to an [Ignition config](https://coreos.github.io/ignition/) for machine consumption.

## Usage

Configure the config transpiler provider (e.g. `providers.tf`).

```tf
provider "ct" {}

terraform {
  required_providers {
    ct = {
      source  = "poseidon/ct"
      version = "0.11.0"
    }
  }
}
```

Define a Butane config for Fedora CoreOS or Flatcar Linux:

```yaml
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

```yaml
variant: flatcar
version: 1.0.0
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

Butane configs are converted to the current (according to this provider) stable Ignition config and merged together. For example, a Butane Config with `version: 1.2.0` would produce an Ignition config with version `v3.3.0`. This relies on Ignition's [forward compatibility](https://github.com/coreos/ignition/blob/main/config/v3_3/config.go#L61).

| terraform-provider-ct | Butane variant | Butane version | Ignition verison |
|-----------------------|----------------|----------------|------------------|
| 0.12.x                | fcos    | 1.0.0, 1.1.0, 1.2.0, 1.3.0, 1.4.0 | 3.3.0 |
| 0.12.x                | flatcar | 1.0.0                             | 3.3.0 |

Before `poseidon/ct` v0.12.0, `ct_config` content could be a Butane Config or a Container Linux Config (CLC). Before `poseidon/ct` v0.10.0, Butane configs contained a `version` that was associated with an Ignition format version. For example, a Butane config with `version: 1.0.0` would produce an Ignition config with version `3.0.0`.

| terraform-provider-ct | CLC to Ignition     | Butane to Ignition    |
|-----------------------|---------------------|--------------------|
| 0.11.x                | Renders 2.3.0       | Butane (1.0, 1.1, 1.2, 1.3, 1.4) -> Ignition 3.3 |
| 0.10.x                | Renders 2.3.0       | Butane (1.0, 1.1, 1.2, 1.3, 1.4) -> Ignition 3.3 |
| 0.9.x                 | Renders 2.3.0       | Butane (1.0, 1.1, 1.2, 1.3, 1.4) -> Ignition (3.0, 3.1, 3.2, 3.2, 3.3)
| 0.8.x                 | Renders 2.3.0       | Butane (1.0, 1.1, 1.2, 1.3) -> Ignition (3.0, 3.1, 3.2, 3.2)
| 0.7.x                 | Renders 2.3.0       | Butane (1.0, 1.1, 1.2) -> Ignition (3.0, 3.1, 3.2) |
| 0.6.x                 | Renders 2.3.0       | Butane (1.0, 1.1) -> Ignition (3.0, 3.1) |
| 0.5.x                 | Renders 2.2.0       | Butane 1.0.0 -> Ignition 3.0.0 |
| 0.4.x                 | Renders 2.2.0       | Butane 1.0.0 -> Ignition 3.0.0 |
| 0.3.x                 | Renders 2.2.0       | NA                 |
| 0.2.x                 | Renders 2.0.0       | NA                 |

## Development

### Binary

To develop the provider plugin locally, build an executable with Go v1.17+.

```
make
```
