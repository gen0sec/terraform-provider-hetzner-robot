package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFailover() *schema.Resource {
	return &schema.Resource{
		Description: "Routes a Hetzner failover IP to an active server (a floating IP you can move between " +
			"servers). Destroying the resource leaves the current routing in place — failover IPs are not deletable.",
		CreateContext: resourceFailoverSet,
		ReadContext:   resourceFailoverRead,
		UpdateContext: resourceFailoverSet,
		DeleteContext: resourceFailoverDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"failover_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The failover IP address",
			},
			"active_server_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP of the server the failover IP currently routes to",
			},
			"netmask":       {Type: schema.TypeString, Computed: true, Description: "Failover netmask"},
			"server_ip":     {Type: schema.TypeString, Computed: true, Description: "IP of the server the failover IP belongs to"},
			"server_number": {Type: schema.TypeInt, Computed: true, Description: "Number of the server the failover IP belongs to"},
		},
	}
}

func resourceFailoverSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	fo, err := c.setFailover(ctx, d.Get("failover_ip").(string), d.Get("active_server_ip").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fo.IP)
	d.Set("active_server_ip", fo.ActiveServerIP)
	d.Set("netmask", fo.Netmask)
	d.Set("server_ip", fo.ServerIP)
	d.Set("server_number", fo.ServerNumber)
	return nil
}

func resourceFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	fo, err := c.getFailover(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("failover_ip", fo.IP)
	d.Set("active_server_ip", fo.ActiveServerIP)
	d.Set("netmask", fo.Netmask)
	d.Set("server_ip", fo.ServerIP)
	d.Set("server_number", fo.ServerNumber)
	return nil
}

func resourceFailoverDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "failover routing left in place",
		Detail:   "Failover IPs are permanent and cannot be deleted; destroying this resource only stops Terraform managing the route. The IP still points to the last active_server_ip.",
	})
	d.SetId("")
	return diags
}
