# Fuze Terraform Provider

The Fuze provider exposes data sources to render [Ignition] [1]
configurations in the human-friendly [Fuze] [2] YAML format into
JSON.  The rendered JSON strings can be used as input to other
Terraform resources, e.g. as user-data for cloud instances.

  [1]: https://github.com/coreos/ignition "Ignition"
  [2]: https://github.com/coreos/fuze     "Fuze"


## Installation

`go get -u github.com/coreos/terraform-provider-fuze`

Update your `.terraformrc` file with the path to the binary:

```hcl
providers {
  fuze = "/path/to/terraform-provider-fuze"
}
```


## Example Usage

```hcl
data "fuze_config" "web" {
  pretty_print = false
  content      = "${file("web.yaml")}"
}

resource "aws_instance" "web" {
    user_data = "${data.fuze_config.web.rendered}"
}
```
