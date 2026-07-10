package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceReset() *schema.Resource {
	return &schema.Resource{
		Description: "Triggers a Hetzner Robot server reset. Action-style resource: it performs the " +
			"reset on create and is a no-op on destroy. Change `triggers` (e.g. to the boot profile) " +
			"to force a new reset — typically chained after hetzner-robot_boot to apply an " +
			"install/rescue profile.",
		CreateContext: resourceResetCreate,
		ReadContext:   resourceResetRead,
		DeleteContext: resourceResetDelete,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the server to reset",
			},
			"reset_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "hw",
				Description:  "Reset type: hw (hardware), sw (Ctrl+Alt+Del), power (power cycle), or man (manual)",
				ValidateFunc: validation.StringInSlice([]string{"hw", "sw", "power", "man"}, false),
			},
			"triggers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Arbitrary key/value pairs that force a new reset when changed (e.g. the boot profile).",
			},
		},
	}
}

func resourceResetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverID := d.Get("server_id").(int)
	if err := c.resetServer(ctx, serverID, d.Get("reset_type").(string)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(serverID))
	return nil
}

// resourceResetRead is a no-op: a reset is a one-shot action with no remote
// state to reconcile.
func resourceResetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// resourceResetDelete only drops the resource from state; it does not undo the reset.
func resourceResetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
