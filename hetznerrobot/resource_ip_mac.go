package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIPMac() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a virtual MAC (vMAC) for an IP address (PUT/DELETE /ip/{ip}/mac).",
		CreateContext: resourceIPMacCreate,
		ReadContext:   resourceIPMacRead,
		DeleteContext: resourceIPMacDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"ip":  {Type: schema.TypeString, Required: true, ForceNew: true, Description: "IP address"},
			"mac": {Type: schema.TypeString, Computed: true, Description: "The generated virtual MAC"},
		},
	}
}

func resourceIPMacCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	mac, err := c.createIPMac(ctx, d.Get("ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(mac.IP)
	d.Set("mac", mac.Mac)
	return nil
}

func resourceIPMacRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	mac, err := c.getIPMac(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("ip", mac.IP)
	d.Set("mac", mac.Mac)
	return nil
}

func resourceIPMacDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	if err := c.deleteIPMac(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
