package hetznerrobot

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRdns() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a reverse DNS (PTR) entry for a Hetzner IP address.",
		CreateContext: resourceRdnsSet,
		ReadContext:   resourceRdnsRead,
		UpdateContext: resourceRdnsSet,
		DeleteContext: resourceRdnsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP address",
			},
			"ptr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PTR record (hostname)",
			},
		},
	}
}

// resourceRdnsSet handles both create and update (POST is an upsert).
func resourceRdnsSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	rdns, err := c.setRdns(ctx, d.Get("ip").(string), d.Get("ptr").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(rdns.IP)
	d.Set("ptr", rdns.PTR)
	return nil
}

func resourceRdnsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	rdns, err := c.getRdns(ctx, d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("ip", rdns.IP)
	d.Set("ptr", rdns.PTR)
	return nil
}

func resourceRdnsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	if err := c.deleteRdns(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
