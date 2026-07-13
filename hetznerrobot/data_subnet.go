package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSubnet() *schema.Resource {
	return &schema.Resource{
		Description: "Reads a subnet's details.",
		ReadContext: dataSubnetRead,
		Schema: map[string]*schema.Schema{
			"subnet_ip":        {Type: schema.TypeString, Required: true, Description: "Subnet network IP"},
			"mask":             {Type: schema.TypeInt, Computed: true},
			"gateway":          {Type: schema.TypeString, Computed: true},
			"server_ip":        {Type: schema.TypeString, Computed: true},
			"server_number":    {Type: schema.TypeInt, Computed: true},
			"failover":         {Type: schema.TypeBool, Computed: true},
			"locked":           {Type: schema.TypeBool, Computed: true},
			"traffic_warnings": {Type: schema.TypeBool, Computed: true},
			"traffic_hourly":   {Type: schema.TypeInt, Computed: true},
			"traffic_daily":    {Type: schema.TypeInt, Computed: true},
			"traffic_monthly":  {Type: schema.TypeInt, Computed: true},
		},
	}
}

func dataSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	sn, err := c.getSubnet(ctx, d.Get("subnet_ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("mask", sn.Mask)
	d.Set("gateway", sn.Gateway)
	d.Set("server_ip", sn.ServerIP)
	d.Set("server_number", sn.ServerNumber)
	d.Set("failover", sn.Failover)
	d.Set("locked", sn.Locked)
	d.Set("traffic_warnings", sn.TrafficWarnings)
	d.Set("traffic_hourly", sn.TrafficHourly)
	d.Set("traffic_daily", sn.TrafficDaily)
	d.Set("traffic_monthly", sn.TrafficMonthly)
	d.SetId(sn.IP)
	return nil
}
