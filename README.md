# Terraform Provider for Container Linux Configs

The `ct` provider provides a `ct_config` data source that parses a [Container Linux Config](https://github.com/coreos/container-linux-config-transpiler/blob/master/doc/configuration.md), validates the content, and renders [Ignition](https://github.com/coreos/ignition). The rendered strings can be used as input to other Terraform resources (e.g. user-data for instances).

## Requirements

* Terraform 0.9.x

## Installation

Add the `terraform-provider-ct` plugin binary on your filesystem.

```
# dev
go get -u github.com/coreos/terraform-provider-ct
```

Register the plugin in `~/.terraformrc`.

```hcl
providers {
  ct = "/$GOPATH/bin/terraform-provider-ct"
}
```

## Usage

```hcl
data "ct_config" "worker" {
  content      = "${file("worker.yaml")}"
  platform     = "ec2"
  pretty_print = false
}

resource "aws_instance" "worker" {
  user_data = "${data.ct_config.worker.rendered}"
}
```

The optional platform can be "azure", "ec2", "gce", or [others](https://github.com/coreos/container-linux-config-transpiler/blob/master/config/platform/platform.go) to perform platform-specific susbstitutions. By default, platform is "" (none, for bare-metal). 

## Development

To develop the provider plugin locally, set up [Go](http://www.golang.org/) 1.8+ and a valid [GOPATH](http://golang.org/doc/code.html#GOPATH). Build the plugin locally.

```sh
make
```

Optionally, install the plugin to `$(GOPATH)/bin`.

```sh
make install
```

### Vendor

Add or update dependencies in glide.yaml and vendor. The [glide](https://github.com/Masterminds/glide) and [glide-vc](https://github.com/sgotti/glide-vc) tools vendor and prune dependencies.

```sh
make vendor
```
