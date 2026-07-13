package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a subnet's traffic-warning configuration.",
		CreateContext: resourceSubnetSet,
		ReadContext:   resourceSubnetRead,
		UpdateContext: resourceSubnetSet,
		DeleteContext: resourceIPNoOpDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"subnet_ip":        {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Subnet network IP"},
			"traffic_warnings": {Type: schema.TypeBool, Optional: true, Default: false},
			"traffic_hourly":   {Type: schema.TypeInt, Optional: true, Default: 0},
			"traffic_daily":    {Type: schema.TypeInt, Optional: true, Default: 0},
			"traffic_monthly":  {Type: schema.TypeInt, Optional: true, Default: 0},
			"mask":             {Type: schema.TypeInt, Computed: true},
			"gateway":          {Type: schema.TypeString, Computed: true},
			"server_ip":        {Type: schema.TypeString, Computed: true},
			"server_number":    {Type: schema.TypeInt, Computed: true},
			"failover":         {Type: schema.TypeBool, Computed: true},
		},
	}
}

func resourceSubnetSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	sn, err := c.setSubnet(ctx, d.Get("subnet_ip").(string), d.Get("traffic_warnings").(bool), d.Get("traffic_hourly").(int), d.Get("traffic_daily").(int), d.Get("traffic_monthly").(int))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(sn.IP)
	setSubnetAttrs(d, sn)
	return nil
}

func resourceSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	sn, err := c.getSubnet(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("subnet_ip", sn.IP)
	setSubnetAttrs(d, sn)
	return nil
}

func setSubnetAttrs(d *schema.ResourceData, sn *HetznerRobotSubnet) {
	d.Set("traffic_warnings", sn.TrafficWarnings)
	d.Set("traffic_hourly", sn.TrafficHourly)
	d.Set("traffic_daily", sn.TrafficDaily)
	d.Set("traffic_monthly", sn.TrafficMonthly)
	d.Set("mask", sn.Mask)
	d.Set("gateway", sn.Gateway)
	d.Set("server_ip", sn.ServerIP)
	d.Set("server_number", sn.ServerNumber)
	d.Set("failover", sn.Failover)
}
