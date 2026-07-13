package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an existing dedicated server's name (POST /server/{number}). Does not create or " +
			"delete servers — destroying the resource only stops Terraform managing the name.",
		CreateContext: resourceServerSet,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerSet,
		DeleteContext: resourceServerNoOpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Server number",
			},
			"server_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server name",
			},
			"server_ip": {Type: schema.TypeString, Computed: true, Description: "Primary IPv4 address"},
		},
	}
}

func resourceServerSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber := d.Get("server_number").(int)

	server, err := c.renameServer(ctx, serverNumber, d.Get("server_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(serverNumber))
	d.Set("server_name", server.ServerName)
	d.Set("server_ip", server.ServerIP)
	return nil
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	server, err := c.getServer(ctx, serverNumber)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("server_number", server.ServerNumber)
	d.Set("server_name", server.ServerName)
	d.Set("server_ip", server.ServerIP)
	return nil
}

func resourceServerNoOpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Servers cannot be deleted via this endpoint; drop from state only.
	d.SetId("")
	return nil
}
