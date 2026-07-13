package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSubnetMac() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a virtual MAC (vMAC) for a subnet (PUT/DELETE /subnet/{ip}/mac).",
		CreateContext: resourceSubnetMacCreate,
		ReadContext:   resourceSubnetMacRead,
		DeleteContext: resourceSubnetMacDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"subnet_ip": {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Subnet network IP"},
			"mac":       {Type: schema.TypeString, Computed: true, Description: "The generated virtual MAC"},
		},
	}
}

func resourceSubnetMacCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	mac, err := c.createSubnetMac(ctx, d.Get("subnet_ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(mac.IP)
	d.Set("mac", mac.Mac)
	return nil
}

func resourceSubnetMacRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	mac, err := c.getSubnetMac(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("subnet_ip", mac.IP)
	d.Set("mac", mac.Mac)
	return nil
}

func resourceSubnetMacDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	if err := c.deleteSubnetMac(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
