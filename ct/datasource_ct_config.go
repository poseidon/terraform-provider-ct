package ct

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"

	ct "github.com/coreos/container-linux-config-transpiler/config"
	ignition "github.com/coreos/ignition/config/v2_1"
	ignitionTypesV2_1 "github.com/coreos/ignition/config/v2_1/types"
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
			"rendered": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "rendered ignition configuration",
			},
		},
	}
}

func dataSourceCTConfigRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderCTConfig(d)
	if err != nil {
		return err
	}

	d.Set("rendered", rendered)
	d.SetId(strconv.Itoa(hashcode.String(rendered)))
	return nil
}

func renderCTConfig(d *schema.ResourceData) (string, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	content := d.Get("content").(string)
	platform := d.Get("platform").(string)
	snippetsIface := d.Get("snippets").([]interface{})
	pretty := d.Get("pretty_print").(bool)

	snippets := make([]string, len(snippetsIface))
	for i := range snippetsIface {
		snippets[i] = snippetsIface[i].(string)
	}

	ign, err := clcToIgnition([]byte(content), platform)
	if err != nil {
		return "", err
	}

	for _, content := range snippets {
		ignext, err := clcToIgnition([]byte(content), platform)
		if err != nil {
			return "", err
		}
		ign = ignition.Append(ign, ignext)
	}

	if pretty {
		ignitionJSON, err := json.MarshalIndent(ign, "", "  ")
		return string(ignitionJSON), err
	}

	ignitionJSON, err := json.Marshal(ign)
	return string(ignitionJSON), err
}

// Parse Container Linux config and convert to Ignition v2.1.0 format.
func clcToIgnition(data []byte, platform string) (ignitionTypesV2_1.Config, error) {
	// parse bytes into a Container Linux Config
	clc, ast, report := ct.Parse([]byte(data))
	if report.IsFatal() {
		return ignitionTypesV2_1.Config{}, fmt.Errorf("error parsing Container Linux Config: %v", report.String())
	}
	// convert Container Linux Config to an Ignition Config
	ign, report := ct.Convert(clc, platform, ast)
	if report.IsFatal() {
		return ignitionTypesV2_1.Config{}, fmt.Errorf("error converting to Ignition: %v", report.String())
	}
	return ign, nil
}
