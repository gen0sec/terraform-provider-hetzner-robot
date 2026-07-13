package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataFailover() *schema.Resource {
	return &schema.Resource{
		Description: "Reads a Hetzner failover IP's current routing.",
		ReadContext: dataFailoverRead,
		Schema: map[string]*schema.Schema{
			"failover_ip":      {Type: schema.TypeString, Required: true, Description: "The failover IP address"},
			"active_server_ip": {Type: schema.TypeString, Computed: true, Description: "IP of the server the failover IP routes to"},
			"netmask":          {Type: schema.TypeString, Computed: true, Description: "Failover netmask"},
			"server_ip":        {Type: schema.TypeString, Computed: true, Description: "IP of the server the failover IP belongs to"},
			"server_number":    {Type: schema.TypeInt, Computed: true, Description: "Number of the server the failover IP belongs to"},
		},
	}
}

func dataFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	fo, err := c.getFailover(ctx, d.Get("failover_ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("active_server_ip", fo.ActiveServerIP)
	d.Set("netmask", fo.Netmask)
	d.Set("server_ip", fo.ServerIP)
	d.Set("server_number", fo.ServerNumber)
	d.SetId(fo.IP)
	return nil
}
