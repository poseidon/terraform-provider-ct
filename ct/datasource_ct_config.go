package ct

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	clct "github.com/coreos/container-linux-config-transpiler/config"
	fcct "github.com/coreos/fcct/config"
	"github.com/coreos/fcct/config/common"
	ignition "github.com/coreos/ignition/config/v2_2"
	ignitionTypesV2_2 "github.com/coreos/ignition/config/v2_2/types"
)

func dataSourceCTConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCTConfigRead,

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

func dataSourceCTConfigRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderConfig(d)
	if err != nil {
		return err
	}

	d.Set("rendered", rendered)
	d.SetId(strconv.Itoa(hashcode.String(rendered)))
	return nil
}

// Render a Fedora CoreOS Config or Container Linux Config as Ignition JSON.
func renderConfig(d *schema.ResourceData) (string, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	content := d.Get("content").(string)
	platform := d.Get("platform").(string)
	snippetsIface := d.Get("snippets").([]interface{})
	pretty := d.Get("pretty_print").(bool)
	strict := d.Get("strict").(bool)

	snippets := make([]string, len(snippetsIface))
	for i := range snippetsIface {
		snippets[i] = snippetsIface[i].(string)
	}

	// Fedora CoreOS Config
	ign, err := fccToIgnition([]byte(content), pretty, strict)
	if err == fcct.ErrNoVariant {
		// consider as Container Linux Config
		ign, err = renderCLC([]byte(content), platform, pretty, strict, snippets)
	}
	return string(ign), err
}

// Translate Fedora CoreOS config to Ignition v3.X.Y
func fccToIgnition(data []byte, pretty, strict bool) ([]byte, error) {
	ign, _, err := fcct.Translate(data, common.TranslateOptions{
		Pretty: pretty,
		Strict: strict,
	})
	return ign, err
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

	if pretty {
		return json.MarshalIndent(ign, "", "  ")
	}
	return json.Marshal(ign)
}

// Parse Container Linux config and convert to Ignition v2.2.0 format.
func clcToIgnition(data []byte, platform string, strict bool) (ignitionTypesV2_2.Config, error) {
	// parse bytes into a Container Linux Config
	clc, ast, report := clct.Parse([]byte(data))

	if strict && len(report.Entries) > 0 {
		return ignitionTypesV2_2.Config{}, fmt.Errorf("error strict parsing Container Linux Config: %v", report.String())
	}

	if report.IsFatal() {
		return ignitionTypesV2_2.Config{}, fmt.Errorf("error parsing Container Linux Config: %v", report.String())
	}
	// convert Container Linux Config to an Ignition Config
	ign, report := clct.Convert(clc, platform, ast)
	if report.IsFatal() {
		return ignitionTypesV2_2.Config{}, fmt.Errorf("error converting to Ignition: %v", report.String())
	}
	return ign, nil
}
