# Container-Linux-Config-Transpiler Terraform Provider

The CT (formerly known as Fuze) provider exposes data sources to render [Ignition] [1]
configurations in the human-friendly [Config-Transpiler] [2] YAML format into
JSON.  The rendered JSON strings can be used as input to other
Terraform resources, e.g. as user-data for cloud instances.

  [1]: https://github.com/coreos/ignition "Ignition"
  [2]: https://github.com/coreos/container-linux-config-transpiler "CT"


## Installation

`go get -u github.com/coreos/terraform-provider-fuze`

Update your `.terraformrc` file with the path to the binary:

```hcl
providers {
  ct = "/$GOPATH/bin/terraform-provider-fuze"
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

## Build

```
make
```

### Dependencies

### Adding Dependencies

After adding a new `import` to the source, use `glide get` to add the dependency to the `glide.yaml` and `glide.lock` files.

```
glide get github.com/$ORG/$PROJ
```

### Updating Dependencies

To update an existing package, edit the `glide.yaml` file to the desired verison (most likely a semver tag or git hash), and run `make revendor`.

```
{{ edit the entry in glide.yaml }}
make revendor
```

If the update was successful, `glide.lock` will have been updated to reflect the changes to `glide.yaml` and the package will have been updated in `vendor`.
