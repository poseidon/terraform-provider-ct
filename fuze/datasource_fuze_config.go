package fuze

import (
	"encoding/json"
	"errors"
	"strconv"

	fuze "github.com/coreos/container-linux-config-transpiler/config"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFuzeConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFuzeConfigRead,

		Schema: map[string]*schema.Schema{
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceFuzeConfigRead(d *schema.ResourceData, meta interface{}) error {
	rendered, err := renderFuzeConfig(d)
	if err != nil {
		return err
	}

	d.Set("rendered", rendered)
	d.SetId(strconv.Itoa(hashcode.String(rendered)))
	return nil
}

func renderFuzeConfig(d *schema.ResourceData) (string, error) {
	pretty := d.Get("pretty_print").(bool)
	config := d.Get("content").(string)

	cfg, rpt := fuze.Parse([]byte(config))
	if rpt.IsFatal() {
		return "", errors.New(rpt.String())
	}

	ignition, rpt := fuze.ConvertAs2_0_0(cfg)
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
