# terraform-provider-ct

Notable changes between releases.

## Latest

## v0.3.2

* Add compatibility with Terraform v0.12. Retain v0.11 compatibility ([#37](https://github.com/poseidon/terraform-provider-ct/pull/37))

## v0.3.1

* Document usage with the Terraform [3rd-party plugin](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) directory ([#35](https://github.com/poseidon/terraform-provider-ct/pull/35))
* Provide pre-compiled binaries from Go v1.11.5
  * Add windows release binaries

## v0.3.0

* Render Ignition Configs at Ingition v2.2.0
* Add `snippets` field for appending Container Linux Configs to `content` ([#22](https://github.com/poseidon/terraform-provider-ct/pull/22))
* Update `ct` to v0.8.0

## v0.2.1

* Add `snippets` field for appending Container Linux Configs to `content` ([#22](https://github.com/poseidon/terraform-provider-ct/pull/22))

## v0.2.0

* Render Ignition Configs at Ignition v2.0.0
* Add `platform` field to the `ct_config` data source
* Add support for platform [dynamic templating](https://coreos.com/os/docs/latest/dynamic-data.html)
* Update to support Terraform v0.9.2+
* Update Container Linux `ct` to v0.3.1

## v0.1.0

Initial release as `tf-provider-fuze`.
