package internal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/clarketm/json"
	"github.com/coreos/go-semver/semver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	butane "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	ignition "github.com/coreos/ignition/v2/config/v3_6"
)

func DatasourceConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceConfigRead,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snippets": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"files_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "allow embedding local files relative to this directory",
			},
			"pretty_print": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"strict": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"use_mapped_version": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, output the ignition version corresponding to the given Butane version, rather than upgrading to the maximum supported version.",
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered ignition configuration",
			},
		},
	}
}

func datasourceConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rendered, err := renderConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rendered", rendered); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode(rendered))
	return diags
}

// Render a Fedora CoreOS Config or Container Linux Config as Ignition JSON.
func renderConfig(d *schema.ResourceData) (string, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	content := d.Get("content").(string)
	pretty := d.Get("pretty_print").(bool)
	filesDir := d.Get("files_dir").(string)
	strict := d.Get("strict").(bool)
	useMappedVersion := d.Get("use_mapped_version").(bool)
	snippetsIface := d.Get("snippets").([]interface{})

	snippets := make([]string, len(snippetsIface))
	for i, v := range snippetsIface {
		if v != nil {
			snippets[i] = v.(string)
		}
	}

	ign, err := butaneToIgnition([]byte(content), pretty, filesDir, strict, snippets, useMappedVersion)
	return string(ign), err
}

// Translate Butane config to Ignition JSON or OpenShift MachineConfig YAML.
func butaneToIgnition(data []byte, pretty bool, filesDir string, strict bool, snippets []string, useMappedVersion bool) ([]byte, error) {
	ignBytes, report, err := butane.TranslateBytes(data, common.TranslateBytesOptions{
		TranslateOptions: common.TranslateOptions{
			FilesDir: filesDir,
		},
		Pretty: pretty,
	})
	// ErrNoVariant indicates data is a CLC, not an FCC
	if err != nil {
		return nil, err
	}
	if strict && len(report.Entries) > 0 {
		return nil, fmt.Errorf("strict parsing error: %v", report.String())
	}

	// OpenShift variant outputs YAML MachineConfig rather than plain ignition JSON
	trimmed := bytes.TrimSpace(ignBytes)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		if len(snippets) > 0 {
			return nil, fmt.Errorf("snippets are not supported for OpenShift configs")
		}
		return ignBytes, nil
	}

	// No snippets and mapped version: return butane output directly (preserves mapped version)
	if len(snippets) == 0 && useMappedVersion {
		return ignBytes, nil
	}

	// merge FCC snippets into main Ignition config
	return mergeFCCSnippets(ignBytes, pretty, filesDir, strict, snippets, useMappedVersion)
}

// Parse Fedora CoreOS Ignition and Butane snippets into Ignition Config.
func mergeFCCSnippets(ignBytes []byte, pretty bool, filesDir string, strict bool, snippets []string, useMappedVersion bool) ([]byte, error) {
	var origVersion string
	var origSemver *semver.Version
	if useMappedVersion {
		var err error
		origVersion, err = ignitionVersion(ignBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to read ignition version: %v", err)
		}
		origSemver, err = semver.NewVersion(origVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid main config ignition version %q: %v", origVersion, err)
		}
	}

	ign, _, err := ignition.ParseCompatibleVersion(ignBytes)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	for _, snippet := range snippets {
		ignextBytes, report, err := butane.TranslateBytes([]byte(snippet), common.TranslateBytesOptions{
			TranslateOptions: common.TranslateOptions{
				FilesDir: filesDir,
			},
			Pretty: pretty,
		})
		if err != nil {
			// For FCC, require snippets be FCCs (don't fall-through to CLC)
			if err == common.ErrNoVariant {
				return nil, fmt.Errorf("Butane snippets require `variant`: %v", err)
			}
			return nil, fmt.Errorf("Butane translate error: %v\n%s", err, report.String())
		}
		if strict && len(report.Entries) > 0 {
			return nil, fmt.Errorf("strict parsing error: %v", report.String())
		}

		if useMappedVersion {
			snippetVersion, err := ignitionVersion(ignextBytes)
			if err != nil {
				return nil, fmt.Errorf("failed to read snippet ignition version: %v", err)
			}
			snippetSemver, err := semver.NewVersion(snippetVersion)
			if err != nil {
				return nil, fmt.Errorf("invalid snippet ignition version %q: %v", snippetVersion, err)
			}
			if origSemver.LessThan(*snippetSemver) {
				return nil, fmt.Errorf("snippet ignition version %s exceeds main config version %s; use a snippet with version <= main config", snippetVersion, origVersion)
			}
		}

		ignext, _, err := ignition.ParseCompatibleVersion(ignextBytes)
		if err != nil {
			return nil, fmt.Errorf("snippet parse error: %v, expect v1.4.0", err)
		}
		ign = ignition.Merge(ign, ignext)
	}

	if useMappedVersion {
		// Restore the original version; ParseCompatibleVersion upgrades to MaxVersion.
		ign.Ignition.Version = origVersion
	}

	return marshalJSON(ign, pretty)
}

// ignitionVersion extracts the ignition version string from a serialized config.
func ignitionVersion(ignBytes []byte) (string, error) {
	var cfg struct {
		Ignition struct {
			Version string `json:"version"`
		} `json:"ignition"`
	}
	if err := json.Unmarshal(ignBytes, &cfg); err != nil {
		return "", err
	}
	return cfg.Ignition.Version, nil
}

func marshalJSON(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}
