// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"reflect"
	"testing"

	"github.com/coreos/container-linux-config-transpiler/config/types"
	ignTypes "github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/config/validate/report"
)

func TestParse(t *testing.T) {
	type in struct {
		data string
	}
	type out struct {
		cfg types.Config
		r   report.Report
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in: in{data: ``},
			out: out{
				cfg: types.Config{},
				r: report.Report{
					Entries: []report.Entry{{
						Kind:    report.EntryWarning,
						Message: "Configuration is empty",
					}},
				},
			},
		},
		{
			in: in{data: `
networkd:
  units:
    - name: bad.blah
      contents: not valid
`},
			out: out{cfg: types.Config{
				Networkd: types.Networkd{
					Units: []types.NetworkdUnit{
						{Name: "bad.blah", Contents: "not valid"},
					},
				},
			}},
		},

		// Config
		{
			in: in{data: `
ignition:
  config:
    append:
      - source: http://example.com/test1
        verification:
          hash:
            function: sha512
            sum: 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
      - source: http://example.com/test2
    replace:
      source: http://example.com/test3
      verification:
        hash:
          function: sha512
          sum: 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
`},
			out: out{cfg: types.Config{
				Ignition: types.Ignition{
					Config: types.IgnitionConfig{
						Append: []types.ConfigReference{
							{
								Source: "http://example.com/test1",
								Verification: types.Verification{
									Hash: types.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							{
								Source: "http://example.com/test2",
							},
						},
						Replace: &types.ConfigReference{
							Source: "http://example.com/test3",
							Verification: types.Verification{
								Hash: types.Hash{
									Function: "sha512",
									Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
								},
							},
						},
					},
				},
			}},
		},

		// Storage
		{
			in: in{data: `
storage:
  disks:
    - device: /dev/sda
      wipe_table: true
      partitions:
        - label: ROOT
          number: 7
          size: 100MB
          start: 50MB
          type_guid: 11111111-1111-1111-1111-111111111111
        - label: DATA
          number: 12
          size: 1GB
          start: 300MB
          type_guid: 00000000-0000-0000-0000-000000000000
        - label: NOTHING
    - device: /dev/sdb
      wipe_table: true
  raid:
    - name: fast
      level: raid0
      devices:
        - /dev/sdc
        - /dev/sdd
    - name: durable
      level: raid1
      devices:
        - /dev/sde
        - /dev/sdf
        - /dev/sdg
      spares: 1
  filesystems:
    - name: filesystem1
      mount:
        device: /dev/disk/by-partlabel/ROOT
        format: btrfs
        create:
          force: true
          options:
            - -L
            - ROOT
    - name: filesystem2
      mount:
        device: /dev/disk/by-partlabel/DATA
        format: ext4
    - name: filesystem3
      path: /sysroot
  files:
    - path: /opt/file1
      filesystem: filesystem1
      contents:
        inline: file1
      mode: 0644
      user:
        id: 500
      group:
        id: 501
    - path: /opt/file2
      filesystem: filesystem1
      contents:
        remote:
          url: http://example.com/file2
          compression: gzip
          verification:
            hash:
              function: sha512
              sum: 00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
      mode: 0644
      user:
        id: 502
      group:
        id: 503
    - path: /opt/file3
      filesystem: filesystem2
      contents:
        remote:
          url: http://example.com/file3
          compression: gzip
      mode: 0400
      user:
        id: 1000
      group:
        id: 1001
    - path: /opt/file4
      filesystem: filesystem2
`},
			out: out{cfg: types.Config{
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device:    "/dev/sda",
							WipeTable: true,
							Partitions: []types.Partition{
								{
									Label:    "ROOT",
									Number:   7,
									Size:     "100MB",
									Start:    "50MB",
									TypeGUID: "11111111-1111-1111-1111-111111111111",
								},
								{
									Label:    "DATA",
									Number:   12,
									Size:     "1GB",
									Start:    "300MB",
									TypeGUID: "00000000-0000-0000-0000-000000000000",
								},
								{
									Label: "NOTHING",
								},
							},
						},
						{
							Device:    "/dev/sdb",
							WipeTable: true,
						},
					},
					Arrays: []types.Raid{
						{
							Name:    "fast",
							Level:   "raid0",
							Devices: []string{"/dev/sdc", "/dev/sdd"},
						},
						{
							Name:    "durable",
							Level:   "raid1",
							Devices: []string{"/dev/sde", "/dev/sdf", "/dev/sdg"},
							Spares:  1,
						},
					},
					Filesystems: []types.Filesystem{
						{
							Name: "filesystem1",
							Mount: &types.Mount{
								Device: "/dev/disk/by-partlabel/ROOT",
								Format: "btrfs",
								Create: &types.Create{
									Force:   true,
									Options: []string{"-L", "ROOT"},
								},
							},
						},
						{
							Name: "filesystem2",
							Mount: &types.Mount{
								Device: "/dev/disk/by-partlabel/DATA",
								Format: "ext4",
							},
						},
						{
							Name: "filesystem3",
							Path: "/sysroot",
						},
					},
					Files: []types.File{
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file1",
							Contents: types.FileContents{
								Inline: "file1",
							},
							Mode:  0644,
							User:  types.FileUser{Id: 500},
							Group: types.FileGroup{Id: 501},
						},
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file2",
							Contents: types.FileContents{
								Remote: types.Remote{
									Url:         "http://example.com/file2",
									Compression: "gzip",
									Verification: types.Verification{
										Hash: types.Hash{
											Function: "sha512",
											Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
										},
									},
								},
							},
							Mode:  0644,
							User:  types.FileUser{Id: 502},
							Group: types.FileGroup{Id: 503},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file3",
							Contents: types.FileContents{
								Remote: types.Remote{
									Url:         "http://example.com/file3",
									Compression: "gzip",
								},
							},
							Mode:  0400,
							User:  types.FileUser{Id: 1000},
							Group: types.FileGroup{Id: 1001},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file4",
							Contents: types.FileContents{
								Inline: "",
							},
						},
					},
				},
			}},
		},

		// systemd
		{
			in: in{data: `
systemd:
  units:
    - name: test1.service
      enable: true
      contents: test1 contents
      dropins:
        - name: conf1.conf
          contents: conf1 contents
        - name: conf2.conf
          contents: conf2 contents
    - name: test2.service
      mask: true
      contents: test2 contents
`},
			out: out{cfg: types.Config{
				Systemd: types.Systemd{
					Units: []types.SystemdUnit{
						{
							Name:     "test1.service",
							Enable:   true,
							Contents: "test1 contents",
							DropIns: []types.SystemdUnitDropIn{
								{
									Name:     "conf1.conf",
									Contents: "conf1 contents",
								},
								{
									Name:     "conf2.conf",
									Contents: "conf2 contents",
								},
							},
						},
						{
							Name:     "test2.service",
							Mask:     true,
							Contents: "test2 contents",
						},
					},
				},
			}},
		},

		// networkd
		{
			in: in{data: `
networkd:
  units:
    - name: empty.netdev
    - name: test.network
      contents: test config
`},
			out: out{cfg: types.Config{
				Networkd: types.Networkd{
					Units: []types.NetworkdUnit{
						{
							Name: "empty.netdev",
						},
						{
							Name:     "test.network",
							Contents: "test config",
						},
					},
				},
			}},
		},

		// passwd
		{
			in: in{data: `
passwd:
  users:
    - name: user 1
      password_hash: password 1
      ssh_authorized_keys:
        - key1
        - key2
    - name: user 2
      password_hash: password 2
      ssh_authorized_keys:
        - key3
        - key4
      create:
        uid: 123
        gecos: gecos
        home_dir: /home/user 2
        no_create_home: true
        primary_group: wheel
        groups:
          - wheel
          - plugdev
        no_user_group: true
        system: true
        no_log_init: true
        shell: /bin/zsh
    - name: user 3
      password_hash: password 3
      ssh_authorized_keys:
        - key5
        - key6
      create: {}
  groups:
    - name: group 1
      gid: 1000
      password_hash: password 1
      system: true
    - name: group 2
      password_hash: password 2
`},
			out: out{cfg: types.Config{
				Passwd: types.Passwd{
					Users: []types.User{
						{
							Name:              "user 1",
							PasswordHash:      "password 1",
							SSHAuthorizedKeys: []string{"key1", "key2"},
						},
						{
							Name:              "user 2",
							PasswordHash:      "password 2",
							SSHAuthorizedKeys: []string{"key3", "key4"},
							Create: &types.UserCreate{
								Uid:          func(i uint) *uint { return &i }(123),
								GECOS:        "gecos",
								Homedir:      "/home/user 2",
								NoCreateHome: true,
								PrimaryGroup: "wheel",
								Groups:       []string{"wheel", "plugdev"},
								NoUserGroup:  true,
								System:       true,
								NoLogInit:    true,
								Shell:        "/bin/zsh",
							},
						},
						{
							Name:              "user 3",
							PasswordHash:      "password 3",
							SSHAuthorizedKeys: []string{"key5", "key6"},
							Create:            &types.UserCreate{},
						},
					},
					Groups: []types.Group{
						{
							Name:         "group 1",
							Gid:          func(i uint) *uint { return &i }(1000),
							PasswordHash: "password 1",
							System:       true,
						},
						{
							Name:         "group 2",
							PasswordHash: "password 2",
						},
					},
				},
			}},
		},
	}

	for i, test := range tests {
		cfg, err := Parse([]byte(test.in.data))
		if !reflect.DeepEqual(err, test.out.r) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.r, err)
		}
		if !reflect.DeepEqual(cfg, test.out.cfg) {
			t.Errorf("#%d: bad config: want %#v, got %#v", i, test.out.cfg, cfg)
		}
	}
}
func TestConvertAs2_0_0(t *testing.T) {
	type in struct {
		cfg types.Config
	}
	type out struct {
		cfg ignTypes.Config
		r   report.Report
	}

	tests := []struct {
		in  in
		out out
	}{
		{
			in:  in{cfg: types.Config{}},
			out: out{cfg: ignTypes.Config{Ignition: ignTypes.Ignition{Version: ignTypes.IgnitionVersion{Major: 2}}}},
		},
		{
			in: in{cfg: types.Config{
				Networkd: types.Networkd{
					Units: []types.NetworkdUnit{
						{Name: "bad.blah", Contents: "not valid"},
					},
				},
			}},
			out: out{r: report.ReportFromError(errors.New("invalid networkd unit extension"), report.EntryError)},
		},

		// Config
		{
			in: in{cfg: types.Config{
				Ignition: types.Ignition{
					Config: types.IgnitionConfig{
						Append: []types.ConfigReference{
							{
								Source: "http://example.com/test1",
								Verification: types.Verification{
									Hash: types.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							{
								Source: "http://example.com/test2",
							},
						},
						Replace: &types.ConfigReference{
							Source: "http://example.com/test3",
							Verification: types.Verification{
								Hash: types.Hash{
									Function: "sha512",
									Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
								},
							},
						},
					},
				},
			}},
			out: out{cfg: ignTypes.Config{
				Ignition: ignTypes.Ignition{
					Version: ignTypes.IgnitionVersion{Major: 2},
					Config: ignTypes.IgnitionConfig{
						Append: []ignTypes.ConfigReference{
							{
								Source: ignTypes.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/test1",
								},
								Verification: ignTypes.Verification{
									Hash: &ignTypes.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							{
								Source: ignTypes.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/test2",
								},
							},
						},
						Replace: &ignTypes.ConfigReference{
							Source: ignTypes.Url{
								Scheme: "http",
								Host:   "example.com",
								Path:   "/test3",
							},
							Verification: ignTypes.Verification{
								Hash: &ignTypes.Hash{
									Function: "sha512",
									Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
								},
							},
						},
					},
				},
			}},
		},

		// Storage
		{
			in: in{cfg: types.Config{
				Storage: types.Storage{
					Disks: []types.Disk{
						{
							Device:    "/dev/sda",
							WipeTable: true,
							Partitions: []types.Partition{
								{
									Label:    "ROOT",
									Number:   7,
									Size:     "100MB",
									Start:    "50MB",
									TypeGUID: "11111111-1111-1111-1111-111111111111",
								},
								{
									Label:    "DATA",
									Number:   12,
									Size:     "1GB",
									Start:    "300MB",
									TypeGUID: "00000000-0000-0000-0000-000000000000",
								},
								{
									Label: "NOTHING",
								},
							},
						},
						{
							Device:    "/dev/sdb",
							WipeTable: true,
						},
					},
					Arrays: []types.Raid{
						{
							Name:    "fast",
							Level:   "raid0",
							Devices: []string{"/dev/sdc", "/dev/sdd"},
						},
						{
							Name:    "durable",
							Level:   "raid1",
							Devices: []string{"/dev/sde", "/dev/sdf", "/dev/sdg"},
							Spares:  1,
						},
					},
					Filesystems: []types.Filesystem{
						{
							Name: "filesystem1",
							Mount: &types.Mount{
								Device: "/dev/disk/by-partlabel/ROOT",
								Format: "btrfs",
								Create: &types.Create{
									Force:   true,
									Options: []string{"-L", "ROOT"},
								},
							},
						},
						{
							Name: "filesystem2",
							Mount: &types.Mount{
								Device: "/dev/disk/by-partlabel/DATA",
								Format: "ext4",
							},
						},
						{
							Name: "filesystem3",
							Path: "/sysroot",
						},
					},
					Files: []types.File{
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file1",
							Contents: types.FileContents{
								Inline: "file1",
							},
							Mode:  0644,
							User:  types.FileUser{Id: 500},
							Group: types.FileGroup{Id: 501},
						},
						{
							Filesystem: "filesystem1",
							Path:       "/opt/file2",
							Contents: types.FileContents{
								Remote: types.Remote{
									Url:         "http://example.com/file2",
									Compression: "gzip",
									Verification: types.Verification{
										Hash: types.Hash{
											Function: "sha512",
											Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
										},
									},
								},
							},
							Mode:  0644,
							User:  types.FileUser{Id: 502},
							Group: types.FileGroup{Id: 503},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file3",
							Contents: types.FileContents{
								Remote: types.Remote{
									Url:         "http://example.com/file3",
									Compression: "gzip",
								},
							},
							Mode:  0400,
							User:  types.FileUser{Id: 1000},
							Group: types.FileGroup{Id: 1001},
						},
						{
							Filesystem: "filesystem2",
							Path:       "/opt/file4",
							Contents: types.FileContents{
								Inline: "",
							},
						},
					},
				},
			}},
			out: out{cfg: ignTypes.Config{
				Ignition: ignTypes.Ignition{Version: ignTypes.IgnitionVersion{Major: 2}},
				Storage: ignTypes.Storage{
					Disks: []ignTypes.Disk{
						{
							Device:    ignTypes.Path("/dev/sda"),
							WipeTable: true,
							Partitions: []ignTypes.Partition{
								{
									Label:    ignTypes.PartitionLabel("ROOT"),
									Number:   7,
									Size:     ignTypes.PartitionDimension(0x32000),
									Start:    ignTypes.PartitionDimension(0x19000),
									TypeGUID: "11111111-1111-1111-1111-111111111111",
								},
								{
									Label:    ignTypes.PartitionLabel("DATA"),
									Number:   12,
									Size:     ignTypes.PartitionDimension(0x200000),
									Start:    ignTypes.PartitionDimension(0x96000),
									TypeGUID: "00000000-0000-0000-0000-000000000000",
								},
								{
									Label: ignTypes.PartitionLabel("NOTHING"),
								},
							},
						},
						{
							Device:    ignTypes.Path("/dev/sdb"),
							WipeTable: true,
						},
					},
					Arrays: []ignTypes.Raid{
						{
							Name:    "fast",
							Level:   "raid0",
							Devices: []ignTypes.Path{ignTypes.Path("/dev/sdc"), ignTypes.Path("/dev/sdd")},
						},
						{
							Name:    "durable",
							Level:   "raid1",
							Devices: []ignTypes.Path{ignTypes.Path("/dev/sde"), ignTypes.Path("/dev/sdf"), ignTypes.Path("/dev/sdg")},
							Spares:  1,
						},
					},
					Filesystems: []ignTypes.Filesystem{
						{
							Name: "filesystem1",
							Mount: &ignTypes.FilesystemMount{
								Device: ignTypes.Path("/dev/disk/by-partlabel/ROOT"),
								Format: ignTypes.FilesystemFormat("btrfs"),
								Create: &ignTypes.FilesystemCreate{
									Force:   true,
									Options: ignTypes.MkfsOptions([]string{"-L", "ROOT"}),
								},
							},
						},
						{
							Name: "filesystem2",
							Mount: &ignTypes.FilesystemMount{
								Device: ignTypes.Path("/dev/disk/by-partlabel/DATA"),
								Format: ignTypes.FilesystemFormat("ext4"),
							},
						},
						{
							Name: "filesystem3",
							Path: func(p ignTypes.Path) *ignTypes.Path { return &p }("/sysroot"),
						},
					},
					Files: []ignTypes.File{
						{
							Filesystem: "filesystem1",
							Path:       ignTypes.Path("/opt/file1"),
							Contents: ignTypes.FileContents{
								Source: ignTypes.Url{
									Scheme: "data",
									Opaque: ",file1",
								},
							},
							Mode:  ignTypes.FileMode(0644),
							User:  ignTypes.FileUser{Id: 500},
							Group: ignTypes.FileGroup{Id: 501},
						},
						{
							Filesystem: "filesystem1",
							Path:       ignTypes.Path("/opt/file2"),
							Contents: ignTypes.FileContents{
								Source: ignTypes.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/file2",
								},
								Compression: "gzip",
								Verification: ignTypes.Verification{
									Hash: &ignTypes.Hash{
										Function: "sha512",
										Sum:      "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
									},
								},
							},
							Mode:  ignTypes.FileMode(0644),
							User:  ignTypes.FileUser{Id: 502},
							Group: ignTypes.FileGroup{Id: 503},
						},
						{
							Filesystem: "filesystem2",
							Path:       ignTypes.Path("/opt/file3"),
							Contents: ignTypes.FileContents{
								Source: ignTypes.Url{
									Scheme: "http",
									Host:   "example.com",
									Path:   "/file3",
								},
								Compression: "gzip",
							},
							Mode:  ignTypes.FileMode(0400),
							User:  ignTypes.FileUser{Id: 1000},
							Group: ignTypes.FileGroup{Id: 1001},
						},
						{
							Filesystem: "filesystem2",
							Path:       ignTypes.Path("/opt/file4"),
							Contents: ignTypes.FileContents{
								Source: ignTypes.Url{
									Scheme: "data",
									Opaque: ",",
								},
							},
						},
					},
				},
			}},
		},

		// systemd
		{
			in: in{cfg: types.Config{
				Systemd: types.Systemd{
					Units: []types.SystemdUnit{
						{
							Name:     "test1.service",
							Enable:   true,
							Contents: "test1 contents",
							DropIns: []types.SystemdUnitDropIn{
								{
									Name:     "conf1.conf",
									Contents: "conf1 contents",
								},
								{
									Name:     "conf2.conf",
									Contents: "conf2 contents",
								},
							},
						},
						{
							Name:     "test2.service",
							Mask:     true,
							Contents: "test2 contents",
						},
					},
				},
			}},
			out: out{cfg: ignTypes.Config{
				Ignition: ignTypes.Ignition{Version: ignTypes.IgnitionVersion{Major: 2}},
				Systemd: ignTypes.Systemd{
					Units: []ignTypes.SystemdUnit{
						{
							Name:     "test1.service",
							Enable:   true,
							Contents: "test1 contents",
							DropIns: []ignTypes.SystemdUnitDropIn{
								{
									Name:     "conf1.conf",
									Contents: "conf1 contents",
								},
								{
									Name:     "conf2.conf",
									Contents: "conf2 contents",
								},
							},
						},
						{
							Name:     "test2.service",
							Mask:     true,
							Contents: "test2 contents",
						},
					},
				},
			}},
		},

		// networkd
		{
			in: in{cfg: types.Config{
				Networkd: types.Networkd{
					Units: []types.NetworkdUnit{
						{
							Name: "empty.netdev",
						},
						{
							Name:     "test.network",
							Contents: "test config",
						},
					},
				},
			}},
			out: out{cfg: ignTypes.Config{
				Ignition: ignTypes.Ignition{Version: ignTypes.IgnitionVersion{Major: 2}},
				Networkd: ignTypes.Networkd{
					Units: []ignTypes.NetworkdUnit{
						{
							Name: "empty.netdev",
						},
						{
							Name:     "test.network",
							Contents: "test config",
						},
					},
				},
			}},
		},

		// passwd
		{
			in: in{cfg: types.Config{
				Passwd: types.Passwd{
					Users: []types.User{
						{
							Name:              "user 1",
							PasswordHash:      "password 1",
							SSHAuthorizedKeys: []string{"key1", "key2"},
						},
						{
							Name:              "user 2",
							PasswordHash:      "password 2",
							SSHAuthorizedKeys: []string{"key3", "key4"},
							Create: &types.UserCreate{
								Uid:          func(i uint) *uint { return &i }(123),
								GECOS:        "gecos",
								Homedir:      "/home/user 2",
								NoCreateHome: true,
								PrimaryGroup: "wheel",
								Groups:       []string{"wheel", "plugdev"},
								NoUserGroup:  true,
								System:       true,
								NoLogInit:    true,
								Shell:        "/bin/zsh",
							},
						},
						{
							Name:              "user 3",
							PasswordHash:      "password 3",
							SSHAuthorizedKeys: []string{"key5", "key6"},
							Create:            &types.UserCreate{},
						},
					},
					Groups: []types.Group{
						{
							Name:         "group 1",
							Gid:          func(i uint) *uint { return &i }(1000),
							PasswordHash: "password 1",
							System:       true,
						},
						{
							Name:         "group 2",
							PasswordHash: "password 2",
						},
					},
				},
			}},
			out: out{cfg: ignTypes.Config{
				Ignition: ignTypes.Ignition{Version: ignTypes.IgnitionVersion{Major: 2}},
				Passwd: ignTypes.Passwd{
					Users: []ignTypes.User{
						{
							Name:              "user 1",
							PasswordHash:      "password 1",
							SSHAuthorizedKeys: []string{"key1", "key2"},
						},
						{
							Name:              "user 2",
							PasswordHash:      "password 2",
							SSHAuthorizedKeys: []string{"key3", "key4"},
							Create: &ignTypes.UserCreate{
								Uid:          func(i uint) *uint { return &i }(123),
								GECOS:        "gecos",
								Homedir:      "/home/user 2",
								NoCreateHome: true,
								PrimaryGroup: "wheel",
								Groups:       []string{"wheel", "plugdev"},
								NoUserGroup:  true,
								System:       true,
								NoLogInit:    true,
								Shell:        "/bin/zsh",
							},
						},
						{
							Name:              "user 3",
							PasswordHash:      "password 3",
							SSHAuthorizedKeys: []string{"key5", "key6"},
							Create:            &ignTypes.UserCreate{},
						},
					},
					Groups: []ignTypes.Group{
						{
							Name:         "group 1",
							Gid:          func(i uint) *uint { return &i }(1000),
							PasswordHash: "password 1",
							System:       true,
						},
						{
							Name:         "group 2",
							PasswordHash: "password 2",
						},
					},
				},
			}},
		},
	}

	for i, test := range tests {
		cfg, r := ConvertAs2_0_0(test.in.cfg)
		if !reflect.DeepEqual(r, test.out.r) {
			t.Errorf("#%d: bad error: want %v, got %v", i, test.out.r, r)
		}
		if !reflect.DeepEqual(cfg, test.out.cfg) {
			t.Errorf("#%d: bad config: want %#v, got %#v", i, test.out.cfg, cfg)
		}
	}
}
