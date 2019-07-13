package ct

import (
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testProviders = map[string]terraform.ResourceProvider{
	"ct": Provider(),
}

const prettyResource = `
data "ct_config" "example" {
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

const prettyExpected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.2.0"
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

const ec2Resource = `
data "ct_config" "ec2" {
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

const ec2Expected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.2.0"
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

const snippetsResource = `
data "ct_config" "combine" {
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
const snippetsExpected = `{
  "ignition": {
    "config": {},
    "security": {
      "tls": {}
    },
    "timeouts": {},
    "version": "2.2.0"
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

const fedoraCoreOSResource = `
data "ct_config" "fedora-coreos" {
  pretty_print = true
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

const fedoraCoreOSExpected = `{
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
  "systemd": {}
}`

func TestRender(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			r.TestStep{
				Config: prettyResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.example", "rendered", prettyExpected),
				),
			},
			r.TestStep{
				Config: ec2Resource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.ec2", "rendered", ec2Expected),
				),
			},
			r.TestStep{
				Config: snippetsResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.combine", "rendered", snippetsExpected),
				),
			},
			r.TestStep{
				Config: fedoraCoreOSResource,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.ct_config.fedora-coreos", "rendered", fedoraCoreOSExpected),
				),
			},
		},
	})
}
