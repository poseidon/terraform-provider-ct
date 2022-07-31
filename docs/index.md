# config-transpiler Provider

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

