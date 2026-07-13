package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataIP() *schema.Resource {
	return &schema.Resource{
		Description: "Reads an IP address's details.",
		ReadContext: dataIPRead,
		Schema: map[string]*schema.Schema{
			"ip":               {Type: schema.TypeString, Required: true, Description: "IP address"},
			"server_ip":        {Type: schema.TypeString, Computed: true},
			"server_number":    {Type: schema.TypeInt, Computed: true},
			"locked":           {Type: schema.TypeBool, Computed: true},
			"separate_mac":     {Type: schema.TypeString, Computed: true},
			"traffic_warnings": {Type: schema.TypeBool, Computed: true},
			"traffic_hourly":   {Type: schema.TypeInt, Computed: true},
			"traffic_daily":    {Type: schema.TypeInt, Computed: true},
			"traffic_monthly":  {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	ip, err := c.getIP(ctx, d.Get("ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("server_ip", ip.ServerIP)
	d.Set("server_number", ip.ServerNumber)
	d.Set("locked", ip.Locked)
	d.Set("separate_mac", ip.SeparateMAC)
	d.Set("traffic_warnings", ip.TrafficWarnings)
	d.Set("traffic_hourly", ip.TrafficHourly)
	d.Set("traffic_daily", ip.TrafficDaily)
	d.Set("traffic_monthly", ip.TrafficMonthly)
	d.SetId(ip.IP)
	return nil
}
