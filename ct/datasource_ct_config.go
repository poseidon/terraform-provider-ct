package ct

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	butane "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	clct "github.com/coreos/container-linux-config-transpiler/config"

	ignition "github.com/coreos/ignition/config/v2_3"
	ignitionTypes "github.com/coreos/ignition/config/v2_3/types"
	ignition33 "github.com/coreos/ignition/v2/config/v3_3"
)

func dataSourceCTConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCTConfigRead,

		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"files_dir": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"platform": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"snippets": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"pretty_print": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"strict": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered ignition configuration",
			},
		},
	}
}

func dataSourceCTConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rendered, err := renderConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("rendered", rendered)
	d.SetId(strconv.Itoa(hashcode.String(rendered)))
	return diags
}

// Render a Fedora CoreOS Config or Container Linux Config as Ignition JSON.
func renderConfig(d *schema.ResourceData) (string, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	content := d.Get("content").(string)
	filesDir := d.Get("files_dir").(string)
	platform := d.Get("platform").(string)
	pretty := d.Get("pretty_print").(bool)
	strict := d.Get("strict").(bool)
	snippetsIface := d.Get("snippets").([]interface{})

	snippets := make([]string, len(snippetsIface))
	for i, v := range snippetsIface {
		if v != nil {
			snippets[i] = v.(string)
		}
	}

	// Butane Config
	ign, err := butaneToIgnition([]byte(content), filesDir, pretty, strict, snippets)
	if err == common.ErrNoVariant {
		// consider as Container Linux Config
		ign, err = renderCLC([]byte(content), platform, pretty, strict, snippets)
	}
	return string(ign), err
}

// Translate Fedora CoreOS config to Ignition v3.X.Y
func butaneToIgnition(data []byte, filesDir string, pretty, strict bool, snippets []string) ([]byte, error) {
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

	// merge FCC snippets into main Ignition config
	return mergeFCCSnippets(ignBytes, filesDir, pretty, strict, snippets)
}

// Parse Fedora CoreOS Ignition and Butane snippets into Ignition Config.
func mergeFCCSnippets(ignBytes []byte, filesDir string, pretty, strict bool, snippets []string) ([]byte, error) {
	ign, _, err := ignition33.ParseCompatibleVersion(ignBytes)
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
			return nil, fmt.Errorf("Butane translate error: %v", err)
		}
		if strict && len(report.Entries) > 0 {
			return nil, fmt.Errorf("strict parsing error: %v", report.String())
		}

		ignext, _, err := ignition33.ParseCompatibleVersion(ignextBytes)
		if err != nil {
			return nil, fmt.Errorf("snippet parse error: %v, expect v1.4.0", err)
		}
		ign = ignition33.Merge(ign, ignext)
	}

	return marshalJSON(ign, pretty)
}

// Translate Container Linux Config as Ignition JSON.
func renderCLC(data []byte, platform string, pretty, strict bool, snippets []string) ([]byte, error) {
	ign, err := clcToIgnition(data, platform, strict)
	if err != nil {
		return nil, err
	}

	for _, snippet := range snippets {
		ignext, err := clcToIgnition([]byte(snippet), platform, strict)
		if err != nil {
			return nil, err
		}
		ign = ignition.Append(ign, ignext)
	}

	return marshalJSON(ign, pretty)
}

// Parse Container Linux config and convert to Ignition v2.2.0 format.
func clcToIgnition(data []byte, platform string, strict bool) (ignitionTypes.Config, error) {
	// parse bytes into a Container Linux Config
	clc, ast, report := clct.Parse([]byte(data))

	if strict && len(report.Entries) > 0 {
		return ignitionTypes.Config{}, fmt.Errorf("error strict parsing Container Linux Config: %v", report.String())
	}

	if report.IsFatal() {
		return ignitionTypes.Config{}, fmt.Errorf("error parsing Container Linux Config: %v", report.String())
	}
	// convert Container Linux Config to an Ignition Config
	ign, report := clct.Convert(clc, platform, ast)
	if report.IsFatal() {
		return ignitionTypes.Config{}, fmt.Errorf("error converting to Ignition: %v", report.String())
	}
	return ign, nil
}

func marshalJSON(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}
