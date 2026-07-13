package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIP() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages an IP address's traffic-warning configuration.",
		CreateContext: resourceIPSet,
		ReadContext:   resourceIPRead,
		UpdateContext: resourceIPSet,
		DeleteContext: resourceIPNoOpDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"ip":               {Type: schema.TypeString, Required: true, ForceNew: true, Description: "IP address"},
			"traffic_warnings": {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable traffic warnings"},
			"traffic_hourly":   {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Hourly traffic warning limit (MB)"},
			"traffic_daily":    {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Daily traffic warning limit (MB)"},
			"traffic_monthly":  {Type: schema.TypeInt, Optional: true, Default: 0, Description: "Monthly traffic warning limit (GB)"},
			"server_ip":        {Type: schema.TypeString, Computed: true, Description: "Server IP the address belongs to"},
			"server_number":    {Type: schema.TypeInt, Computed: true, Description: "Server number"},
			"separate_mac":     {Type: schema.TypeString, Computed: true, Description: "Separate MAC, if set"},
		},
	}
}

func resourceIPSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	ip, err := c.setIP(ctx, d.Get("ip").(string), d.Get("traffic_warnings").(bool), d.Get("traffic_hourly").(int), d.Get("traffic_daily").(int), d.Get("traffic_monthly").(int))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ip.IP)
	setIPAttrs(d, ip)
	return nil
}

func resourceIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	ip, err := c.getIP(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ip", ip.IP)
	setIPAttrs(d, ip)
	return nil
}

func resourceIPNoOpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func setIPAttrs(d *schema.ResourceData, ip *HetznerRobotIP) {
	d.Set("traffic_warnings", ip.TrafficWarnings)
	d.Set("traffic_hourly", ip.TrafficHourly)
	d.Set("traffic_daily", ip.TrafficDaily)
	d.Set("traffic_monthly", ip.TrafficMonthly)
	d.Set("server_ip", ip.ServerIP)
	d.Set("server_number", ip.ServerNumber)
	d.Set("separate_mac", ip.SeparateMAC)
}
