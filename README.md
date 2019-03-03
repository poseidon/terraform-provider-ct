# terraform-provider-ct

`terraform-provider-ct` allows defining a `ct_config` [Container Linux Config](https://github.com/coreos/container-linux-config-transpiler/blob/master/doc/configuration.md) resource to validate the content and render an [Ignition](https://github.com/coreos/ignition) document. Rendered Ignition may be used as input to other Terraform resources (e.g. instance user-data).

## Requirements

* Terraform v0.11+ [installed](https://www.terraform.io/downloads.html)

## Install

Add the `terraform-provider-ct` plugin binary for your system to the Terraform 3rd-party [plugin directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) `~/.terraform.d/plugins`.

```sh
VERSION=v0.3.0
wget https://github.com/coreos/terraform-provider-ct/releases/download/$VERSION/terraform-provider-ct-$VERSION-linux-amd64.tar.gz
tar xzf terraform-provider-ct-$VERSION-linux-amd64.tar.gz
mv terraform-provider-ct-$VERSION-linux-amd64/terraform-provider-ct ~/.terraform.d/plugins/terraform-provider-ct_$VERSION
```

Terraform plugin binary names are versioned to allow for migrations of managed infrastructure.

```
$ tree ~/.terraform.d/
/home/user/.terraform.d/
└── plugins
    ├── terraform-provider-ct_v0.2.1
    └── terraform-provider-ct_v0.3.0
```

## Usage

Configure the ct provider in a `providers.tf` file.

```hcl
provider "ct" {
  version = "0.3.0"
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

Declare a `ct_config` resource in Terraform.

```hcl
data "ct_config" "worker" {
  content      = "${file("worker.yaml")}"
  platform     = "ec2"
  pretty_print = false

  snippets = [
    "${file("units.yaml")}",
    "${file("storage.yaml")}",
  ]
}

resource "aws_instance" "worker" {
  user_data = "${data.ct_config.worker.rendered}"
}
```

Set the `content` to the contents of a Container Linux Config that should be validated and rendered as Ignition. Optionally, use the `snippets` field to append a list of Container Linux Config snippets. Use `platform` if [platform-specific](https://github.com/coreos/container-linux-config-transpiler/blob/master/config/platform/platform.go) susbstitution is desired.

### Ignition schema output

Each minor version of `terraform-provider-ct` is tightly coupled with a minor version of the Ignition schema. Ignition transparently handles old Ignition schema versions, so this isn't normally an issue.

Upgrading between versions for existing managed clusters **may not be safe**. Check your usage to determine whether a `user_data` change to an instance would re-provision a machine (important machines may  be configured to ignore `user_data` changes).

| terraform-provider-ct | Ignition schema version |
|-----------------------|-------------------------|
| 0.2.x                 | 2.0                     |
| 0.3.x                 | 2.2                     |

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
