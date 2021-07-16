# terraform-provider-ct

Notable changes between releases.

## Latest

## v0.9.0

* Add Butane Config v1.4.0 support ([#100](https://github.com/poseidon/terraform-provider-ct/pull/100))
  * Accept Butane v1.4.0 config/snippets and render Ignition v3.3.0
* Rename Fedora CoreOS Configs to Butane Configs
* Remove Go module vendoring
* Remove tarball release format

## v0.8.0

* Migrate to Terraform Plugin SDK v2.3.0 ([#75](https://github.com/poseidon/terraform-provider-ct/pull/75))
* Add Fedora CoreOS Config v1.3.0 support ([#76](https://github.com/poseidon/terraform-provider-ct/pull/76))

## v0.7.1

* Fix possible empty rendered Ignition ([#72](https://github.com/poseidon/terraform-provider-ct/pull/72))
  * Fix regression in rendering Fedora CoreOS v1.1.0 Configs with `snippets` and `pretty_print = false`
* Remove Terraform v0.12.x instructions

## v0.7.0

* Add Fedora CoreOS Config v1.2.0 support ([#71](https://github.com/poseidon/terraform-provider-ct/pull/71))
  * Accept FCC v1.2.0 and output Ignition v3.2.0

## v0.6.1

* Fix zip archive artifacts for Darwin and Windows ([#67](https://github.com/poseidon/terraform-provider-ct/pull/67))
* Add Linux ARM64 release artifacts ([#66](https://github.com/poseidon/terraform-provider-ct/pull/66))

## v0.6.0

* Add Fedora CoreOS Config v1.1.0 support ([#63](https://github.com/poseidon/terraform-provider-ct/pull/63))
  * Accept FCC v1.1.0 and output Ignition v3.1.0
  * Continue to support FCC v1.0.0 and output Ignition v3.0.0
  * Support merging FCC snippets into v1.0.0 or v1.1.0 FCC content
  * Note: Version skew among snippets and content is not supported
* Change Container Linux Config to render Ignition v2.3.0 ([#60](https://github.com/poseidon/terraform-provider-ct/pull/60))
* Add zip archive format with signed checksum

## v0.5.1

* Allow empty string snippets ([#61](https://github.com/poseidon/terraform-provider-ct/pull/61))

## v0.5.0

* Add support for Fedora CoreOS `snippets` ([#58](https://github.com/poseidon/terraform-provider-ct/pull/58))
* Migrate to the Terraform Plugin SDK ([#56](https://github.com/poseidon/terraform-provider-ct/pull/56))
* Upgrade fcct from v0.1.0 to v0.4.0 ([#57](https://github.com/poseidon/terraform-provider-ct/pull/57))

## v0.4.0

* Support Fedora CoreOS Config `content` ([#52](https://github.com/poseidon/terraform-provider-ct/pull/52))
  * Render Container Linux Config `content` in Ignition v2.2 format
  * Render Fedora CoreOS `content` in Ignition v3.x format (determined by their own [version](https://github.com/coreos/fcct/blob/master/docs/configuration-v1_0.md))
* Add `strict` field for strict validation (defaults to false) ([#53](https://github.com/poseidon/terraform-provider-ct/pull/53))

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
