package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	butane "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	ignition_rel "github.com/coreos/ignition/v2/config/v3_5"
	types_rel "github.com/coreos/ignition/v2/config/v3_5/types"
	ignition_exp "github.com/coreos/ignition/v2/config/v3_6_experimental"
	types_exp "github.com/coreos/ignition/v2/config/v3_6_experimental/types"
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
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered ignition configuration",
			},
			"experimental": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "use experimental version in ignition configuration",
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
	snippetsIface := d.Get("snippets").([]interface{})
	experimental := d.Get("experimental").(bool)

	snippets := make([]string, len(snippetsIface))
	for i, v := range snippetsIface {
		if v != nil {
			snippets[i] = v.(string)
		}
	}

	// Butane Config
	ign, err := butaneToIgnition([]byte(content), pretty, filesDir, strict, snippets, experimental)
	return string(ign), err
}

// Translate Fedora CoreOS config to Ignition v3.X.Y
func butaneToIgnition(data []byte, pretty bool, filesDir string, strict bool, snippets []string, experimental bool) ([]byte, error) {
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
	return mergeFCCSnippets(ignBytes, pretty, filesDir, strict, snippets, experimental)
}

// Parse Fedora CoreOS Ignition and Butane snippets into Ignition Config.
func mergeFCCSnippets(ignBytes []byte, pretty bool, filesDir string, strict bool, snippets []string, experimental bool) ([]byte, error) {
	ign, _, err := ignParseCompatibleVersion(ignBytes, experimental)
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

		ignext, _, err := ignParseCompatibleVersion(ignextBytes, experimental)
		if err != nil {
			return nil, fmt.Errorf("snippet parse error: %v, expect v1.4.0", err)
		}
		ign = ignMerge(ign, ignext, experimental)
	}

	return marshalJSON(ign, pretty)
}

func ignMerge(parent, child interface{}, experimental bool) any {
	if experimental {
		parent_exp := parent.(types_exp.Config)
		child_exp := child.(types_exp.Config)
		return ignition_exp.Merge(parent_exp, child_exp)
	}
	parent_rel := parent.(types_rel.Config)
	child_rel := child.(types_rel.Config)
	return ignition_rel.Merge(parent_rel, child_rel)
}

func ignParseCompatibleVersion(raw []byte, experimental bool) (any, any, error) {
	if experimental {
		return ignition_exp.ParseCompatibleVersion(raw)
	}
	return ignition_rel.ParseCompatibleVersion(raw)
}

func marshalJSON(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}
