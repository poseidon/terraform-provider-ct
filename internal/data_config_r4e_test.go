package internal

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// RHEL for Edge (r4e) variant, v1.1.0 -> ignition 3.4.0

const r4eV11Resource = `
data "ct_config" "r4e" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: r4e
version: 1.1.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const r4eV11WithSnippets = `
data "ct_config" "r4e-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: r4e
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
variant: r4e
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const r4eV11WithSnippetsPrettyFalse = `
data "ct_config" "r4e-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: r4e
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
variant: r4e
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_R4E_v1_1(t *testing.T) {
	ign := r4eIgnVersion["1.1.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: r4eV11Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.r4e", "rendered", ignExpectNoLuks(ign)),
			},
			{
				Config: r4eV11WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.r4e-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: r4eV11WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.r4e-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// RHEL for Edge (r4e) variant, v1.0.0 -> ignition 3.3.0

const r4eV10Resource = `
data "ct_config" "r4e" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: r4e
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const r4eV10WithSnippets = `
data "ct_config" "r4e-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: r4e
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
variant: r4e
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const r4eV10WithSnippetsPrettyFalse = `
data "ct_config" "r4e-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: r4e
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
variant: r4e
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_R4E_v1_0(t *testing.T) {
	ign := r4eIgnVersion["1.0.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: r4eV10Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.r4e", "rendered", ignExpectNoLuks(ign)),
			},
			{
				Config: r4eV10WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.r4e-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: r4eV10WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.r4e-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}
