# config-transpiler Provider

`terraform-provider-ct` allows Terraform to validate a [Container Linux Config](https://github.com/coreos/container-linux-config-transpiler/blob/master/doc/configuration.md) or [Fedora CoreOS Config](https://github.com/coreos/fcct/blob/master/docs/configuration-v1_1.md) and transpile it as [Ignition](https://github.com/coreos/ignition) for machine consumption.

## Usage

Configure the config transpiler provider (e.g. `providers.tf`).

```hcl
provider "ct" {
  version = "0.6.0"
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

