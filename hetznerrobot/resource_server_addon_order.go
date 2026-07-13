package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerAddonOrder() *schema.Resource {
	return &schema.Resource{
		Description: "Orders an add-on (e.g. extra IP, subnet) for an existing server (POST /order/server_addon/transaction). " +
			"By default test = true (validate only, no charge). Set test = false to place a real, billable order.",
		CreateContext: resourceServerAddonOrderCreate,
		ReadContext:   resourceServerAddonOrderRead,
		DeleteContext: resourceServerOrderDelete, // no-op-with-warning
		Schema: map[string]*schema.Schema{
			"server_number": {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Server number to add the add-on to"},
			"product_id":    {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Add-on product ID (see the hetzner-robot_server_addons data source)"},
			"reason":        {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Justification required by some add-ons (e.g. an additional/primary IPv4). Required by the API for those; ignored otherwise."},
			"test": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "If true (default), validate the order only — no charge is made. Set to false to place a real, billable order.",
			},

			"transaction_id": {Type: schema.TypeString, Computed: true, Description: "Order transaction ID"},
			"status":         {Type: schema.TypeString, Computed: true, Description: "Order status"},
		},
	}
}

func resourceServerAddonOrderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(HetznerRobotClient)

	serverNumber := d.Get("server_number").(int)
	test := d.Get("test").(bool)

	tx, err := c.createServerAddonOrder(ctx, serverNumber, d.Get("product_id").(string), d.Get("reason").(string), test)
	if err != nil {
		return diag.FromErr(err)
	}

	if tx.ID != "" {
		d.SetId(tx.ID)
	} else {
		d.SetId("addon-order-" + strconv.Itoa(serverNumber))
	}
	d.Set("transaction_id", tx.ID)
	d.Set("status", tx.Status)

	if test {
		diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: "test add-on order only", Detail: "test = true validated the add-on order without charging. Set test = false to place a real order."})
	} else {
		diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: "real, billable add-on order placed", Detail: "A billable add-on order was placed. Destroying this resource does NOT remove the add-on."})
	}
	return diags
}

func resourceServerAddonOrderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	tx, err := c.getServerAddonOrder(ctx, d.Id())
	if err != nil {
		return nil
	}
	d.Set("transaction_id", tx.ID)
	d.Set("status", tx.Status)
	return nil
}
