# Container-Linux-Config-Transpiler Terraform Provider

The CT (formerly known as Fuze) provider exposes data sources to render [Ignition] [1]
configurations in the human-friendly [Config-Transpiler] [2] YAML format into
JSON.  The rendered JSON strings can be used as input to other
Terraform resources, e.g. as user-data for cloud instances.

  [1]: https://github.com/coreos/ignition "Ignition"
  [2]: https://github.com/coreos/container-linux-config-transpiler "CT"

## Requirements

* Terraform 0.9.x

## Installation

`go get -u github.com/coreos/terraform-provider-ct`

Update your `.terraformrc` file with the path to the binary:

```hcl
providers {
  ct = "/$GOPATH/bin/terraform-provider-ct"
}
```

## Example Usage

```hcl
data "ct_config" "web" {
  pretty_print = false
  content      = "${file("web.yaml")}"
}

resource "aws_instance" "web" {
  user_data = "${data.ct_config.web.rendered}"
}
```

## Development

To develop the provider plugin locally, you'll need [Go](http://www.golang.org/) 1.8+ installed and a [GOPATH](http://golang.org/doc/code.html#GOPATH) setup. Build the plugin locally.

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
