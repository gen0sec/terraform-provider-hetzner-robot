package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataRdns() *schema.Resource {
	return &schema.Resource{
		Description: "Looks up the reverse DNS (PTR) record for an IP address.",
		ReadContext: dataRdnsRead,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address",
			},
			"ptr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PTR record (hostname)",
			},
		},
	}
}

func dataRdnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	rdns, err := c.getRdns(ctx, d.Get("ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ptr", rdns.PTR)
	d.SetId(rdns.IP)
	return nil
}
