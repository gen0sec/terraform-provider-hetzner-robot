package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWOL() *schema.Resource {
	return &schema.Resource{
		Description: "Sends a Wake-on-LAN packet to a server. Action-style resource: it sends the " +
			"packet on create and is a no-op on destroy. Change `triggers` to send another packet.",
		CreateContext: resourceWOLCreate,
		ReadContext:   resourceWOLRead,
		DeleteContext: resourceWOLDelete,
		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Number of the server to wake",
			},
			"triggers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Arbitrary key/value pairs that force a new WOL packet when changed.",
			},
		},
	}
}

func resourceWOLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber := d.Get("server_number").(int)
	if err := c.sendWOL(ctx, serverNumber); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(serverNumber))
	return nil
}

func resourceWOLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceWOLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
