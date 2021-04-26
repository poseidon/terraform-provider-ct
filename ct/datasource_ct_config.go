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
	ignition30 "github.com/coreos/ignition/v2/config/v3_0"
	ignition30Types "github.com/coreos/ignition/v2/config/v3_0/types"
	ignition31 "github.com/coreos/ignition/v2/config/v3_1"
	ignition31Types "github.com/coreos/ignition/v2/config/v3_1/types"
	ignition32 "github.com/coreos/ignition/v2/config/v3_2"
	ignition32Types "github.com/coreos/ignition/v2/config/v3_2/types"
)

func dataSourceCTConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCTConfigRead,

		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

	// Fedora CoreOS Config
	ign, err := fccToIgnition([]byte(content), pretty, strict, snippets)
	if err == common.ErrNoVariant {
		// consider as Container Linux Config
		ign, err = renderCLC([]byte(content), platform, pretty, strict, snippets)
	}
	return string(ign), err
}

// Translate Fedora CoreOS config to Ignition v3.X.Y
func fccToIgnition(data []byte, pretty, strict bool, snippets []string) ([]byte, error) {
	ignBytes, _, err := butane.TranslateBytes(data, common.TranslateBytesOptions{
		Pretty: pretty,
		Strict: strict,
	})
	// ErrNoVariant indicates data is a CLC, not an FCC
	if err != nil {
		return nil, err
	}

	if len(snippets) == 0 {
		return ignBytes, nil
	}

	// merge FCC snippets into main Ignition config
	return mergeFCCSnippets(ignBytes, pretty, strict, snippets)
}

// Manually parse main Fedora CoreOS Config's Ignition using fallback Ignition
// versions. Then translate and parse FCC snippets as the chosen Ignition
// version to merge.
// version
// Upstream might later handle: https://github.com/coreos/butane/issues/118
// Note: This means snippets version must match the main config version.
func mergeFCCSnippets(ignBytes []byte, pretty, strict bool, snippets []string) ([]byte, error) {
	ign, _, err := ignition32.Parse(ignBytes)
	if err == nil {
		// FCC config v1.2.0
		ign, err = mergeFCC12(ign, snippets, pretty, strict)
		if err != nil {
			return nil, fmt.Errorf("FCC v1.2.0 merge error: %v", err)
		}
		return marshalJSON(ign, pretty)
	}

	ign31, _, err := ignition31.Parse(ignBytes)
	if err == nil {
		// FCC config v1.1.0
		ign31, err = mergeFCC11(ign31, snippets, pretty, strict)
		if err != nil {
			return nil, fmt.Errorf("FCC v1.1.0 merge error: %v", err)
		}
		return marshalJSON(ign31, pretty)
	}

	var ign30 ignition30Types.Config
	ign30, _, err = ignition30.Parse(ignBytes)
	if err != nil {
		return nil, fmt.Errorf("FCC v1.0.0 parse error: %v", err)
	}
	// FCC config v1.0.0
	ign30, err = mergeFCCV10(ign30, snippets, pretty, strict)
	if err != nil {
		return nil, fmt.Errorf("FCC v1.0.0 merge error: %v", err)
	}
	return marshalJSON(ign30, pretty)
}

// merge FCC v1.2.0 snippets
func mergeFCC12(ign ignition32Types.Config, snippets []string, pretty, strict bool) (ignition32Types.Config, error) {
	for _, snippet := range snippets {
		ignextBytes, _, err := butane.TranslateBytes([]byte(snippet), common.TranslateBytesOptions{
			Pretty: pretty,
			Strict: strict,
		})
		if err != nil {
			// For FCC, require snippets be FCCs (don't fall-through to CLC)
			if err == common.ErrNoVariant {
				return ign, fmt.Errorf("Fedora CoreOS snippets require `variant`: %v", err)
			}
			return ign, fmt.Errorf("snippet v1.2.0 translate error: %v", err)
		}
		ignext, _, err := ignition32.Parse(ignextBytes)
		if err != nil {
			return ign, fmt.Errorf("snippet parse error: %v, expect v1.2.0", err)
		}
		ign = ignition32.Merge(ign, ignext)
	}
	return ign, nil
}

// merge FCC v1.1.0 snippets
func mergeFCC11(ign ignition31Types.Config, snippets []string, pretty, strict bool) (ignition31Types.Config, error) {
	for _, snippet := range snippets {
		ignextBytes, _, err := butane.TranslateBytes([]byte(snippet), common.TranslateBytesOptions{
			Pretty: pretty,
			Strict: strict,
		})
		if err != nil {
			// For FCC, require snippets be FCCs (don't fall-through to CLC)
			if err == common.ErrNoVariant {
				return ign, fmt.Errorf("Fedora CoreOS snippets require `variant`: %v", err)
			}
			return ign, fmt.Errorf("snippet v1.1.0 translate error: %v", err)
		}
		ignext, _, err := ignition31.Parse(ignextBytes)
		if err != nil {
			return ign, fmt.Errorf("snippet parse error: %v, expect v1.1.0", err)
		}
		ign = ignition31.Merge(ign, ignext)
	}
	return ign, nil
}

// merge FCC v1.0.0 snippets
func mergeFCCV10(ign ignition30Types.Config, snippets []string, pretty, strict bool) (ignition30Types.Config, error) {
	for _, snippet := range snippets {
		ignextBytes, _, err := butane.TranslateBytes([]byte(snippet), common.TranslateBytesOptions{
			Pretty: pretty,
			Strict: strict,
		})
		if err != nil {
			// For FCC, require snippets be FCCs (don't fall-through to CLC)
			if err == common.ErrNoVariant {
				return ign, fmt.Errorf("Fedora CoreOS snippets require `variant`: %v", err)
			}
			return ign, fmt.Errorf("snippet v1.0.0 translate error: %v", err)
		}
		ignext, _, err := ignition30.Parse(ignextBytes)
		if err != nil {
			return ign, fmt.Errorf("snippet parse error: %v, expect v1.0.0", err)
		}
		ign = ignition30.Merge(ign, ignext)
	}
	return ign, nil
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
