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
EOT
}
`

const prettyExpected = `{
  "ignition": {
    "version": "2.0.0",
    "config": {}
  },
  "storage": {
    "filesystems": [
      {
        "name": "rootfs",
        "mount": {
          "device": "/dev/disk/by-label/ROOT",
          "format": "ext4"
        }
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
        "mode": 420,
        "user": {},
        "group": {}
      }
    ]
  },
  "systemd": {},
  "networkd": {},
  "passwd": {}
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
		},
	})
}
