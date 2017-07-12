package ct

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"

	ct "github.com/coreos/container-linux-config-transpiler/config"
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
	config := d.Get("content").(string)
	platform := d.Get("platform").(string)
	pretty := d.Get("pretty_print").(bool)

	// parse bytes int a Container Linux Config
	cfg, rpt := ct.Parse([]byte(config))
	if rpt.IsFatal() {
		return "", errors.New(rpt.String())
	}

	// convert Container Linux Config to an Ignition Config
	ignition, rpt := ct.ConvertAs2_0(cfg, platform)
	if rpt.IsFatal() {
		return "", errors.New(rpt.String())
	}

	if pretty {
		ignitionJSON, err := json.MarshalIndent(&ignition, "", "  ")
		return string(ignitionJSON), err
	}

	ignitionJSON, err := json.Marshal(&ignition)
	return string(ignitionJSON), err
}
