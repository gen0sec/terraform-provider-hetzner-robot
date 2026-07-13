package hetznerrobot

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataTraffic() *schema.Resource {
	return &schema.Resource{
		Description: "Queries traffic statistics for an IP address over a period.",
		ReadContext: dataTrafficRead,
		Schema: map[string]*schema.Schema{
			"ip": {Type: schema.TypeString, Required: true, Description: "IP address to query"},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Aggregation type: day, month or year",
				ValidateFunc: validation.StringInSlice([]string{"day", "month", "year"}, false),
			},
			"from": {Type: schema.TypeString, Required: true, Description: "Start of the period (e.g. 2024-01-01 or 2024-01-01T10)"},
			"to":   {Type: schema.TypeString, Required: true, Description: "End of the period"},
			"in":   {Type: schema.TypeString, Computed: true, Description: "Incoming traffic (GB)"},
			"out":  {Type: schema.TypeString, Computed: true, Description: "Outgoing traffic (GB)"},
			"sum":  {Type: schema.TypeString, Computed: true, Description: "Total traffic (GB)"},
		},
	}
}

func dataTrafficRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	ip := d.Get("ip").(string)
	t, err := c.getTraffic(ctx, ip, d.Get("type").(string), d.Get("from").(string), d.Get("to").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("in", t.In)
	d.Set("out", t.Out)
	d.Set("sum", t.Sum)
	d.SetId(fmt.Sprintf("%s-%s-%s-%s", ip, d.Get("type").(string), d.Get("from").(string), d.Get("to").(string)))
	return nil
}
