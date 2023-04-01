package internal

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Fedora CoreOS variant, v1.5.0

const fedoraCoreOSV15Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
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
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV15Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignitionV34Expected),
				),
			},
			{
				Config: fedoraCoreOSV15WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV15WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.4.0

const fedoraCoreOSV14Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
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

const ignitionV34Expected = `{
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

const fedoraCoreOSV14WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
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

const ignitionV34WithSnippetsExpected = `{
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

const fedoraCoreOSV14WithSnippetsPrettyFalse = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = false
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

const ignitionV34WithSnippetsPrettyFalseExpected = `{"ignition":{"config":{"replace":{"verification":{}}},"proxy":{},"security":{"tls":{}},"timeouts":{},"version":"3.4.0"},"kernelArguments":{},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"storage":{},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`

func TestButaneConfig_FCOSv1_4(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV14Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignitionV34Expected),
				),
			},
			{
				Config: fedoraCoreOSV14WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV14WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.3.0

const fedoraCoreOSV13Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
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
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV13Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignitionV34Expected),
				),
			},
			{
				Config: fedoraCoreOSV13WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV13WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.2.0

const fedoraCoreOSV12Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
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
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV12Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignitionV34Expected),
				),
			},
			{
				Config: fedoraCoreOSV12WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV12WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.1.0

const fedoraCoreOSV11Resource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
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

// Butane v1.2 added storage.luks, which we exercise
const ignitionV34BeforeButaneV12 = `{
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
  "systemd": {}
}`

const fedoraCoreOSV11WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
  pretty_print = true
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
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV11Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignitionV34BeforeButaneV12),
				),
			},
			{
				Config: fedoraCoreOSV11WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV11WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

// Fedora CoreOS variant, v1.0.0

const fedoraCoreOSV10Resource = `
data "ct_config" "fedora-coreos" {
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

const fedoraCoreOSV10WithSnippets = `
data "ct_config" "fedora-coreos-snippets" {
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
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSV10Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", ignitionV34BeforeButaneV12),
				),
			},
			{
				Config: fedoraCoreOSV10WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsExpected),
				),
			},
			{
				Config: fedoraCoreOSV10WithSnippetsPrettyFalse,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", ignitionV34WithSnippetsPrettyFalseExpected),
				),
			},
		},
	})
}

const fedoraCoreOSMixSnippetBehind = `
data "ct_config" "fedora-coreos-mix-versions" {
  pretty_print = true
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

const ignitionV34MixExpected = `{
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

func TestFedoraCoreOSMix_SnippetBehind(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: fedoraCoreOSMixSnippetBehind,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-mix-versions", "rendered", ignitionV34MixExpected),
				),
			},
		},
	})
}

const fedoraCoreOSMixSnippetAhead = `
data "ct_config" "fedora-coreos-mix-versions" {
  pretty_print = true
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
				Config: fedoraCoreOSMixSnippetAhead,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-mix-versions", "rendered", ignitionV34MixExpected),
				),
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
