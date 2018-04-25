# terraform-provider-ct

Notable changes between releases.

## v0.3.0 (2018-04-25)

* Render Ignition Configs at Ingition 2.2.0
* Update `ct` to v0.8.0
* Add `snippets` field, a list of Container Linux Configs additively appended to the Container Linux Config in `content`

## v0.2.0 (2017-08-03)

* Add `platform` field to the `ct_config` data source
* Add support for platform [dynamic templating](https://coreos.com/os/docs/latest/dynamic-data.html)
* Update to support Terraform v0.9.2+
* Update Container Linux `ct` to v0.3.1

## v0.1.0 (2016-11-11)

Initial release as `tf-provider-fuze`.
