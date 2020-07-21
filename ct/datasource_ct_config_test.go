package ct

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testProviders = map[string]terraform.ResourceProvider{
	"ct": Provider(),
}

const containerLinuxResource = `
data "ct_config" "container-linux" {
  pretty_print = true
  content = <<EOT
---
storage:
  filesystems:
    - name: "rootfs"
      mount:
        device: "/dev/disk/by-label/ROOT"
        format: "ext4"

  files:
    - path: "/etc/motd"
      filesystem: "rootfs"
      mode: 0644
      contents:
        inline: |
          Hello World!
  directories:
    - path: /foo/bar
      filesystem: "rootfs"
      mode: 0644
      user:
        name: root
      group:
        name: root
EOT
}
`

const containerLinuxExpected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.3.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {
    "directories": [
      {
        "filesystem": "rootfs",
        "group": {
          "name": "root"
        },
        "path": "/foo/bar",
        "user": {
          "name": "root"
        },
        "mode": 420
      }
    ],
    "files": [
      {
        "filesystem": "rootfs",
        "path": "/etc/motd",
        "contents": {
          "source": "data:,Hello%20World!%0A",
          "verification": {}
        },
        "mode": 420
      }
    ],
    "filesystems": [
      {
        "mount": {
          "device": "/dev/disk/by-label/ROOT",
          "format": "ext4"
        },
        "name": "rootfs"
      }
    ]
  },
  "systemd": {}
}`

const containerLinuxPlatformResource = `
data "ct_config" "container-linux-platform" {
  pretty_print = true
	platform = "ec2"
  content  = <<EOT
---
etcd:
  advertise_client_urls:       http://{PUBLIC_IPV4}:2379
  initial_advertise_peer_urls: http://{PRIVATE_IPV4}:2380
  listen_client_urls:          http://0.0.0.0:2379
  listen_peer_urls:            http://{PRIVATE_IPV4}:2380
  discovery:                   https://discovery.etcd.io/<token>
EOT
}
`

const containerLinuxPlatformExpected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.3.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {},
  "systemd": {
    "units": [
      {
        "dropins": [
          {
            "contents": "[Unit]\nRequires=coreos-metadata.service\nAfter=coreos-metadata.service\n\n[Service]\nEnvironmentFile=/run/metadata/coreos\nExecStart=\nExecStart=/usr/lib/coreos/etcd-wrapper $ETCD_OPTS \\\n  --listen-peer-urls=\"http://${COREOS_EC2_IPV4_LOCAL}:2380\" \\\n  --listen-client-urls=\"http://0.0.0.0:2379\" \\\n  --initial-advertise-peer-urls=\"http://${COREOS_EC2_IPV4_LOCAL}:2380\" \\\n  --advertise-client-urls=\"http://${COREOS_EC2_IPV4_PUBLIC}:2379\" \\\n  --discovery=\"https://discovery.etcd.io/\u003ctoken\u003e\"",
            "name": "20-clct-etcd-member.conf"
          }
        ],
        "enable": true,
        "name": "etcd-member.service"
      }
    ]
  }
}`

const containerLinuxSnippetsResource = `
data "ct_config" "container-linux-snippets" {
	pretty_print = true
	content = <<EOT
---
storage:
  filesystems:
    - name: "rootfs"
      mount:
        device: "/dev/disk/by-label/ROOT"
        format: "ext4"
EOT
	snippets = [
<<EOT
---
systemd:
  units:
    - name: docker.service
      enable: true
EOT
	]
}
`
const containerLinuxSnippetsExpected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.3.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {
    "filesystems": [
      {
        "mount": {
          "device": "/dev/disk/by-label/ROOT",
          "format": "ext4"
        },
        "name": "rootfs"
      }
    ]
  },
  "systemd": {
    "units": [
      {
        "enable": true,
        "name": "docker.service"
      }
    ]
  }
}`

func TestContainerLinuxConfig(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: containerLinuxResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.container-linux", "rendered", containerLinuxExpected),
				),
			},
			r.TestStep{
				Config: containerLinuxPlatformResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.container-linux-platform", "rendered", containerLinuxPlatformExpected),
				),
			},
			r.TestStep{
				Config: containerLinuxSnippetsResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.container-linux-snippets", "rendered", containerLinuxSnippetsExpected),
				),
			},
		},
	})
}

const fedoraCoreOSResource = `
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

const fedoraCoreOSExpected = `{
  "ignition": {
    "version": "3.1.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  }
}`

const fedoraCoreOSWithSnippets = `
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

const fedoraCoreOSWithSnippetsExpected = `{
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
    "version": "3.1.0"
  },
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

func TestFedoraCoreOSConfigV11(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: fedoraCoreOSResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", fedoraCoreOSExpected),
				),
			},
			r.TestStep{
				Config: fedoraCoreOSWithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", fedoraCoreOSWithSnippetsExpected),
				),
			},
		},
	})
}

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

const fedoraCoreOSV10Expected = `{
  "ignition": {
    "version": "3.0.0"
  },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "key"
        ]
      }
    ]
  }
}`

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

const fedoraCoreOSV10WithSnippetsExpected = `{
  "ignition": {
    "config": {
      "replace": {
        "source": null,
        "verification": {}
      }
    },
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "3.0.0"
  },
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

func TestFedoraCoreOSConfigV10(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: fedoraCoreOSV10Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", fedoraCoreOSV10Expected),
				),
			},
			r.TestStep{
				Config: fedoraCoreOSV10WithSnippets,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos-snippets", "rendered", fedoraCoreOSV10WithSnippetsExpected),
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
version: 1.0.0
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
			r.TestStep{
				Config:      fedoraCoreOSMixSnippetBehind,
				ExpectError: regexp.MustCompile("^.*FCC v1.1.0 merge error: snippet parse error: unsupported config version, expect v1.1.0"),
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
version: 1.1.0
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
			r.TestStep{
				Config:      fedoraCoreOSMixSnippetAhead,
				ExpectError: regexp.MustCompile("^.*FCC v1.0.0 merge error: snippet parse error: unsupported config version, expect v1.0.0"),
			},
		},
	})
}

const invalidResource = `
data "ct_config" "container-linux-strict" {
  content = "foo"
  strict = true
  some_invalid_field = "strict-mode-will-reject"
}
`

func TestInvalidResource(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config:      invalidResource,
				ExpectError: regexp.MustCompile("^.*An argument named \"some_invalid_field\" is not expected here"),
			},
		},
	})
}

// forbidden in strict mode
const emptySnippet = `
data "ct_config" "empty-snippet" {
	pretty_print = true
	content = ""
	snippets = [""]
}
`

const emptySnippetExpected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.3.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {},
  "systemd": {}
}`

func TestAllowEmptySnippet(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: emptySnippet,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.empty-snippet", "rendered", emptySnippetExpected),
				),
			},
		},
	})
}
