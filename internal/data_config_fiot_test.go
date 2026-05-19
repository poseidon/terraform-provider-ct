package internal

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Fedora IoT (fiot) variant, v1.0.0 -> ignition 3.4.0

const fiotV10Resource = `
data "ct_config" "fiot" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fiot
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

const fiotV10WithSnippets = `
data "ct_config" "fiot-snippets" {
  pretty_print = true
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fiot
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
variant: fiot
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const fiotV10WithSnippetsPrettyFalse = `
data "ct_config" "fiot-snippets" {
  pretty_print = false
  use_mapped_version = true
  strict = true
  content = <<EOT
---
variant: fiot
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
variant: fiot
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_Fiot_v1_0(t *testing.T) {
	ign := fiotIgnVersion["1.0.0"]
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fiotV10Resource,
				Check:  r.TestCheckResourceAttr("data.ct_config.fiot", "rendered", ignExpectNoLuks(ign)),
			},
			{
				Config: fiotV10WithSnippets,
				Check:  r.TestCheckResourceAttr("data.ct_config.fiot-snippets", "rendered", ignExpectWithSnippets(ign)),
			},
			{
				Config: fiotV10WithSnippetsPrettyFalse,
				Check:  r.TestCheckResourceAttr("data.ct_config.fiot-snippets", "rendered", ignExpectWithSnippetsCompact(ign)),
			},
		},
	})
}

// Default upgrade behavior tests (use_mapped_version = false)

const fiotDefaultUpgrade = `
data "ct_config" "fiot-default" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: fiot
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - key
EOT
}
`

func TestButaneConfig_Fiot_DefaultUpgrade(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fiotDefaultUpgrade,
				Check:  r.TestCheckResourceAttr("data.ct_config.fiot-default", "rendered", ignExpectNoLuks("3.6.0")),
			},
		},
	})
}
