package internal

import "fmt"

// fcosIgnVersion maps fcos butane spec versions to the ignition spec version they produce.
var fcosIgnVersion = map[string]string{
	"1.0.0": "3.0.0",
	"1.1.0": "3.1.0",
	"1.2.0": "3.2.0",
	"1.3.0": "3.2.0",
	"1.4.0": "3.3.0",
	"1.5.0": "3.4.0",
	"1.6.0": "3.5.0",
	"1.7.0": "3.6.0",
}

// flatcarIgnVersion maps flatcar butane spec versions to ignition spec versions.
var flatcarIgnVersion = map[string]string{
	"1.0.0": "3.3.0",
	"1.1.0": "3.4.0",
}

// r4eIgnVersion maps r4e butane spec versions to ignition spec versions.
var r4eIgnVersion = map[string]string{
	"1.0.0": "3.3.0",
	"1.1.0": "3.4.0",
}

// fiotIgnVersion maps fiot butane spec versions to ignition spec versions.
var fiotIgnVersion = map[string]string{
	"1.0.0": "3.4.0",
}

// ignExpectWithLuks returns the expected pretty-printed ignition JSON for a
// config with a single luks device and no snippets.
func ignExpectWithLuks(ignVersion string) string {
	return fmt.Sprintf(`{
  "ignition": {
    "version": %q
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
  "storage": {
    "luks": [
      {
        "device": "/dev/vdb",
        "name": "data"
      }
    ]
  }
}`, ignVersion)
}

// ignExpectNoLuks returns the expected pretty-printed ignition JSON for a
// passwd-only config with no snippets.
func ignExpectNoLuks(ignVersion string) string {
	return fmt.Sprintf(`{
  "ignition": {
    "version": %q
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
}`, ignVersion)
}

// ignExpectWithSnippets returns the expected pretty-printed ignition JSON after
// merging a systemd snippet into a passwd-only main config.
func ignExpectWithSnippets(ignVersion string) string {
	return fmt.Sprintf(`{
  "ignition": {
    "version": %q
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
  "systemd": {
    "units": [
      {
        "enabled": true,
        "name": "docker.service"
      }
    ]
  }
}`, ignVersion)
}

// ignExpectWithSnippetsCompact returns the expected compact ignition JSON after
// merging a systemd snippet into a passwd-only main config.
func ignExpectWithSnippetsCompact(ignVersion string) string {
	return fmt.Sprintf(`{"ignition":{"version":%q},"passwd":{"users":[{"name":"core","sshAuthorizedKeys":["key"]}]},"systemd":{"units":[{"enabled":true,"name":"docker.service"}]}}`, ignVersion)
}
