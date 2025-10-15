# terraform-provider-ct

Notable changes between releases.

## Latest

## v0.14.0

* Update Butane from v0.24.0 to v0.25.1 ([#223](https://github.com/poseidon/terraform-provider-ct/pull/223), [#219](https://github.com/poseidon/terraform-provider-ct/pull/219), [#212](https://github.com/poseidon/terraform-provider-ct/pull/212), [#199](https://github.com/poseidon/terraform-provider-ct/pull/199))
* Update Ignition from v2.22.0 to v2.24.0 ([#226](https://github.com/poseidon/terraform-provider-ct/pull/226), [#220](https://github.com/poseidon/terraform-provider-ct/pull/220), [#215](https://github.com/poseidon/terraform-provider-ct/pull/215))
* Update Terraform Plugin SDK from v2.36.0 to v2.38.1 ([#222](https://github.com/poseidon/terraform-provider-ct/pull/222), [#221](https://github.com/poseidon/terraform-provider-ct/pull/221))
* Add experimental `files_dir` support to allow embedding local files relative to directory ([5d67e66](https://github.com/poseidon/terraform-provider-ct/commit/5d67e66))
* Improve error messages by including Butane translate reports ([8782f05](https://github.com/poseidon/terraform-provider-ct/commit/8782f05))
* Update Go from v1.19 to v1.24.0 with toolchain v1.25.3

## v0.13.0

* Update the target stable Ignition spec version to v3.4.0 ([#156](https://github.com/poseidon/terraform-provider-ct/pull/156))
  * Parse Butane Configs to Ignition v3.4.0 ([#159](https://github.com/poseidon/terraform-provider-ct/pull/159))
  * Add support for `fcos` [v1.5.0](https://coreos.github.io/butane/config-fcos-v1_5/) Butane Configs
  * Add support for `flatcar` [v1.1.0](https://coreos.github.io/butane/config-flatcar-v1_1/) Butane Configs
* Remove deprecated `platform` field
* Move implementation to an `internal` package ([#157](https://github.com/poseidon/terraform-provider-ct/pull/157))

## v0.12.0

* Remove support for Container Linux Configs ([#132](https://github.com/poseidon/terraform-provider-ct/pull/132))
  * Butane Configs support `fcos` and `flatcar` variants
  * Focus on converting Butane Configs (with different variants) to Ignition
  * Flatcar Linux now supports Ignition v3.3.0
* Remove unused `github.com/coroes/ignition` (v1) dependencies
* Deprecate the `platform` field, it's no longer used

## v0.11.0

* Update coreos/butane from v0.14.0 to v0.15.0 ([#126](https://github.com/poseidon/terraform-provider-ct/pull/126))
  * Add `flatcar` Butane Config variant with spec version 1.0.0 (generates Ignition v3.3.0)
* Deprecate Container Linux Configs (please migrate to Butane Configs)
* Update Go version to v1.18

## v0.10.0

* Change how older (< 1.4) Butane Configs are parsed to Ignition ([#116](https://github.com/poseidon/terraform-provider-ct/pull/116))
  * Parse Ignition bytes to the forward compatible Ignition version ([docs](https://github.com/poseidon/terraform-provider-ct#versions))
  * Parse v1.3 Butane Configs to Ignition v3.3
  * Parse v1.2 Butane Configs to Ignition v3.3
  * Parse v1.1 Butane Configs to Ignition v3.3
  * Parse v1.0 Butane Configs to Ignition v3.3
* Add support for verison skew among Butane Config snippets
  * Butane Config and snippets will always convert to the current Ignition version

## v0.9.2

* Update butane, ignition, and Terraform SDK modules

## v0.9.1

* Update Go version to v1.16+
* Add `darwin-arm64` release target

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
