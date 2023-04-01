package ct

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Flatcar variant, v1.1.0

const flatcarV11Resource = `
data "ct_config" "flatcar" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: flatcar
version: 1.1.0
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

const flatcarV11WithSnippets = `
data "ct_config" "flatcar-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: flatcar
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
variant: flatcar
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const flatcarV11WithSnippetsPrettyFalse = `
data "ct_config" "flatcar-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: flatcar
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
variant: flatcar
version: 1.1.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

func TestButaneConfig_Flatcar_v1_1(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: flatcarV11Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.flatcar", "rendered", flatcarExpected),
				),
			},
			{
				Config: flatcarV11WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.flatcar-snippets", "rendered", flatcarWithSnippetsExpected),
				),
			},
			{
				Config: flatcarV11WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.flatcar-snippets", "rendered", flatcarWithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Flatcar variant, v1.0.0

const flatcarV10Resource = `
data "ct_config" "flatcar" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: flatcar
version: 1.0.0
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

const flatcarExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.4.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {
    "luks": [
      {
        "clevis": {
          "custom": {}
        },
        "device": "/dev/vdb",
        "keyFile": {
          "verification": {}
        },
        "name": "data"
      }
    ]
  },
  "systemd": {}
}`

const flatcarV10WithSnippets = `
data "ct_config" "flatcar-snippets" {
  pretty_print = true
  strict = true
  content = <<EOT
---
variant: flatcar
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
variant: flatcar
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const flatcarWithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "verification": {}
      }
    },
    "proxy": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.4.0"
  },
  "kernelArguments": {},
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  },
  "storage": {},
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`

const flatcarV10WithSnippetsPrettyFalse = `
data "ct_config" "flatcar-snippets" {
  pretty_print = false
  strict = true
  content = <<EOT
---
variant: flatcar
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
variant: flatcar
version: 1.0.0
systemd:
  units:
    - name: docker.service
      enabled: true
EOT
	]
}
`

const flatcarWithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

func TestButaneConfig_Flatcar_v1_0(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: flatcarV10Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.flatcar", "rendered", flatcarExpected),
				),
			},
			{
				Config: flatcarV10WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.flatcar-snippets", "rendered", flatcarWithSnippetsExpected),
				),
			},
			{
				Config: flatcarV10WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.flatcar-snippets", "rendered", flatcarWithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}
