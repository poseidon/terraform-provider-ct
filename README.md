# terraform-provider-ct
[![GoDoc](https://pkg.go.dev/badge/github.com/poseidon/terraform-provider-ct.svg)](https://pkg.go.dev/github.com/poseidon/terraform-provider-ct)
[![Workflow](https://github.com/poseidon/terraform-provider-ct/actions/workflows/test.yaml/badge.svg)](https://github.com/poseidon/terraform-provider-ct/actions/workflows/test.yaml?query=branch%3Amain)
![Downloads](https://img.shields.io/github/downloads/poseidon/terraform-provider-ct/total)
[![Sponsors](https://img.shields.io/github/sponsors/poseidon?logo=github)](https://github.com/sponsors/poseidon)
[![Mastodon](https://img.shields.io/badge/follow-news-6364ff?logo=mastodon)](https://fosstodon.org/@poseidon)

`terraform-provider-ct` allows Terraform to validate a [Butane config](https://coreos.github.io/butane/specs/) and transpile to an [Ignition config](https://coreos.github.io/ignition/) for machine consumption.

## Usage

Configure the config transpiler provider (e.g. `providers.tf`).

```tf
provider "ct" {}

terraform {
  required_providers {
    ct = {
      source  = "poseidon/ct"
      version = "0.14.0"
    }
  }
}
```

Define a Butane config for Fedora CoreOS or Flatcar Linux:

```yaml
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

```yaml
variant: flatcar
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ssh-key foo
```

Define a `ct_config` data source with strict validation.

```tf
data "ct_config" "worker" {
  content      = file("worker.yaml")
  strict       = true
  pretty_print = false

  snippets = [
    file("units.yaml"),
    file("storage.yaml"),
  ]
}
```

Files referenced in the Butane configuration must be placed in the `files_dir` directory. This is equivalent to the butane argument `--files-dir`.

The file is located at `./config/ssh-port.conf`.

```tf
data "ct_config" "worker" {
  content      = file("worker.yaml")
  strict       = true
  pretty_print = false
  files_dir    = "./config"
}
```

```yaml
variant: fcos
version: 1.5.0
storage:
  files:
    - path: /etc/ssh/sshd_config.d/99-ssh-port.conf
      contents:
        local: ssh-port.conf
```

Optionally, template the `content`.

```tf
data "ct_config" "worker" {
  content = templatefile("worker.yaml", {
    ssh_authorized_key = "ssh-ed25519 AAAA...",
  })
  strict       = true
}
```

Render the `ct_config` as Ignition for use by machine instances.

```tf
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

Butane configs are converted to the current (according to this provider) stable Ignition config and merged together. For example, `poseidon/ct` `v0.12.0` would convert a Butane Config with `variant: fcos` and `version: 1.2.0` to an Ignition config with version `v3.3.0`. This relies on Ignition's [forward compatibility](https://github.com/coreos/ignition/blob/main/config/v3_3/config.go#L61).

| poseidon/ct           | Butane variant | Butane version | Ignition verison |
|-----------------------|----------------|----------------|------------------|
| 0.14.x                | fcos    | 1.0.0, 1.1.0, 1.2.0, 1.3.0, 1.4.0, 1.5.0 | 3.4.0 |
| 0.14.x                | flatcar | 1.0.0, 1.1.0                      | 3.4.0 |
| 0.13.x                | fcos    | 1.0.0, 1.1.0, 1.2.0, 1.3.0, 1.4.0, 1.5.0 | 3.4.0 |
| 0.13.x                | flatcar | 1.0.0, 1.1.0                      | 3.4.0 |
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

To develop the provider plugin locally, build an executable with Go v1.18+.

```
make
```
