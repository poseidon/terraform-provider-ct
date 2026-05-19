package internal

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Fedora CoreOS variant, v1.7.0

const fedoraCoreOSV17Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.7.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV17WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.7.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.7.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV17WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.7.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.7.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_7(t *testing.T) {
	ign := fcosIgnVersion["1.7.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV17Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectWithLuks(ign)),
			},
			{
				Config: fedoraCoreOSV17WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV17WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.6.0

const fedoraCoreOSV16Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.6.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV16WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.6.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.6.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV16WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.6.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.6.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_6(t *testing.T) {
	ign := fcosIgnVersion["1.6.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV16Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectWithLuks(ign)),
			},
			{
				Config: fedoraCoreOSV16WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV16WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.5.0

const fedoraCoreOSV15Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.5.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV15WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.5.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV15WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.5.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.5.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_5(t *testing.T) {
	ign := fcosIgnVersion["1.5.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV15Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectWithLuks(ign)),
			},
			{
				Config: fedoraCoreOSV15WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV15WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.4.0

const fedoraCoreOSV14Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV14WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.4.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV14WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.4.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_4(t *testing.T) {
	ign := fcosIgnVersion["1.4.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV14Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectWithLuks(ign)),
			},
			{
				Config: fedoraCoreOSV14WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV14WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.3.0

const fedoraCoreOSV13Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.3.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV13WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.3.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.3.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV13WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.3.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.3.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_3(t *testing.T) {
	ign := fcosIgnVersion["1.3.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV13Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectWithLuks(ign)),
			},
			{
				Config: fedoraCoreOSV13WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV13WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.2.0

const fedoraCoreOSV12Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
storage:
  luks:
    - name: data
      device: /dev/vdb
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV12WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.2.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV12WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.2.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_2(t *testing.T) {
	ign := fcosIgnVersion["1.2.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV12Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectWithLuks(ign)),
			},
			{
				Config: fedoraCoreOSV12WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV12WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.1.0

const fedoraCoreOSV11Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV11WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV11WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_1(t *testing.T) {
	ign := fcosIgnVersion["1.1.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV11Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectNoLuks(ign)),
			},
			{
				Config: fedoraCoreOSV11WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV11WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Fedora CoreOS variant, v1.0.0

const fedoraCoreOSV10Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fedoraCoreOSV10WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fedoraCoreOSV10WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_FCOSv1_0(t *testing.T) {
	ign := fcosIgnVersion["1.0.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV10Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignExpectNoLuks(ign)),
			},
			{
				Config: fedoraCoreOSV10WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fedoraCoreOSV10WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Mixed-version snippet tests: all configurations are upgraded to max version.

const fedoraCoreOSMixSnippetBehind = `
data "ct_config" "fedora-coreos-mix-versions" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.4.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.2.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestFedoraCoreOSMix_SnippetBehind(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSMixSnippetBehind,
				Check:  r.TestCheckResourceAttr("data.ct_config.fedora-coreos-mix-versions", "rendered", ignExpectWithSnippets("3.3.0")),
			},
		},
	})
}

const fedoraCoreOSMixSnippetAhead = `
data "ct_config" "fedora-coreos-mix-versions" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.2.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
	snippets = [
<<EOT
---
variant: fcos
version: 1.4.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestFedoraCoreOSMixVersions_SnippetAhead(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config:      fedoraCoreOSMixSnippetAhead,
				ExpectError: regexp.MustCompile(`snippet ignition version 3\.3\.0 exceeds main config version 3\.2\.0`),
			},
		},
	})
}

// Default upgrade behavior tests (use_mapped_version = false)

const fedoraCoreOSDefaultUpgrade = `
data "ct_config" "default-upgrade" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

func TestFedoraCoreOS_DefaultUpgrade(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSDefaultUpgrade,
				Check:  r.TestCheckResourceAttr("data.ct_config.default-upgrade", "rendered", ignExpectNoLuks("3.6.0")),
			},
		},
	})
}

const invalidResource = `
data "ct_config" "invalid" {
  content = "foo"
  strict = true
  some_invalid_field = "strict-mode-will-reject"
}
`

func TestInvalidResource(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config:      invalidResource,
				ExpectError: regexp.MustCompile("An argument named \"some_invalid_field\" is not expected here"),
			},
		},
	})
}
