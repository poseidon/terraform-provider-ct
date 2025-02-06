package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	butane "github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
	ignition "github.com/coreos/ignition/v2/config/v3_4"
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
	strict := d.Get("strict").(bool)
	snippetsIface := d.Get("snippets").([]interface{})

	snippets := make([]string, len(snippetsIface))
	for i, v := range snippetsIface {
		if v != nil {
			snippets[i] = v.(string)
		}
	}

	// Butane Config
	ign, err := butaneToIgnition([]byte(content), pretty, strict, snippets)
	return string(ign), err
}

// Translate Fedora CoreOS config to Ignition v3.X.Y
func butaneToIgnition(data []byte, pretty, strict bool, snippets []string) ([]byte, error) {
	ignBytes, report, err := butane.TranslateBytes(data, common.TranslateBytesOptions{
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
	return mergeFCCSnippets(ignBytes, pretty, strict, snippets)
}

// Parse Fedora CoreOS Ignition and Butane snippets into Ignition Config.
func mergeFCCSnippets(ignBytes []byte, pretty, strict bool, snippets []string) ([]byte, error) {
	ign, _, err := ignition.ParseCompatibleVersion(ignBytes)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	for _, snippet := range snippets {
		ignextBytes, report, err := butane.TranslateBytes([]byte(snippet), common.TranslateBytesOptions{
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

		ignext, _, err := ignition.ParseCompatibleVersion(ignextBytes)
		if err != nil {
			return nil, fmt.Errorf("snippet parse error: %v, expect v1.4.0", err)
		}
		ign = ignition.Merge(ign, ignext)
	}

	return marshalJSON(ign, pretty)
}

func marshalJSON(v interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}
