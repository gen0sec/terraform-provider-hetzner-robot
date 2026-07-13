package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerOrder() *schema.Resource {
	return &schema.Resource{
		Description: "Orders a Hetzner dedicated server (POST /order/server/transaction). " +
			"By default test = true, which only validates the order without provisioning or " +
			"charging. Set test = false to place a REAL, billable order. Destroying this " +
			"resource does NOT cancel the server — cancel it separately via Hetzner Robot.",
		CreateContext: resourceServerOrderCreate,
		ReadContext:   resourceServerOrderRead,
		DeleteContext: resourceServerOrderDelete,
		Schema: map[string]*schema.Schema{
			"product_id":      {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Product ID to order (see the hetzner-robot_server_products data source)"},
			"location":        {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Datacenter location (e.g. FSN1)"},
			"dist":            {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Distribution/OS to preinstall"},
			"lang":            {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Installation language"},
			"authorized_keys": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "SSH key fingerprints authorized on the ordered server"},
			"password":        {Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true, Description: "Root password (alternative to authorized_keys)"},
			"addons":          {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Product add-on IDs"},
			"primary_ipv4": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Convenience flag: when true, adds the \"primary_ipv4\" add-on so the server is ordered with a public IPv4 (servers are IPv4-less by default). Equivalent to including \"primary_ipv4\" in addons.",
			},
			"test": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "If true (default), validate the order only — no server is provisioned and no charge is made. Set to false to place a real, billable order.",
			},

			"transaction_id": {Type: schema.TypeString, Computed: true, Description: "Order transaction ID"},
			"status":         {Type: schema.TypeString, Computed: true, Description: "Order status"},
			"server_number":  {Type: schema.TypeInt, Computed: true, Description: "Assigned server number (once provisioned)"},
			"server_ip":      {Type: schema.TypeString, Computed: true, Description: "Assigned primary IP (once provisioned)"},
		},
	}
}

func serverOrderStringList(d *schema.ResourceData, key string) []string {
	out := []string{}
	if v, ok := d.GetOk(key); ok {
		for _, item := range v.([]interface{}) {
			out = append(out, item.(string))
		}
	}
	return out
}

// serverOrderAddons returns the explicit addon IDs, plus "primary_ipv4" when the
// primary_ipv4 convenience flag is set (deduplicated). Shared by the server and
// server-market order resources, both of which expose a primary_ipv4 flag.
func serverOrderAddons(d *schema.ResourceData) []string {
	addons := serverOrderStringList(d, "addons")
	if d.Get("primary_ipv4").(bool) {
		for _, a := range addons {
			if a == "primary_ipv4" {
				return addons // already present
			}
		}
		addons = append(addons, "primary_ipv4")
	}
	return addons
}

func resourceServerOrderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(HetznerRobotClient)

	req := HetznerRobotServerOrderRequest{
		ProductID:      d.Get("product_id").(string),
		Location:       d.Get("location").(string),
		Dist:           d.Get("dist").(string),
		Lang:           d.Get("lang").(string),
		AuthorizedKeys: serverOrderStringList(d, "authorized_keys"),
		Password:       d.Get("password").(string),
		Addons:         serverOrderAddons(d),
		Test:           d.Get("test").(bool),
	}

	tx, err := c.createServerOrder(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if tx.ID != "" {
		d.SetId(tx.ID)
	} else {
		// test/validation responses may not carry a transaction id
		d.SetId("order-" + req.ProductID)
	}
	d.Set("transaction_id", tx.ID)
	d.Set("status", tx.Status)
	d.Set("server_number", tx.ServerNumber)
	d.Set("server_ip", tx.ServerIP)

	if req.Test {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "test order only (no server provisioned)",
			Detail:   "test = true validated the order without provisioning or charging. Set test = false to place a real order.",
		})
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "real, billable server order placed",
			Detail:   "A billable dedicated-server order was placed. Provisioning is asynchronous; server_number/server_ip populate once the order is ready. Destroying this resource does NOT cancel the server.",
		})
	}
	return diags
}

func resourceServerOrderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	tx, err := c.getServerOrder(ctx, d.Id())
	if err != nil {
		// test orders / expired transactions may not be retrievable; keep state.
		return nil
	}
	d.Set("transaction_id", tx.ID)
	d.Set("status", tx.Status)
	d.Set("server_number", tx.ServerNumber)
	d.Set("server_ip", tx.ServerIP)
	return nil
}

func resourceServerOrderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	// Orders cannot be cancelled through the ordering API; just drop from state.
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "order not cancelled",
		Detail:   "Removing this resource only removes it from Terraform state; a provisioned server is NOT cancelled. Cancel it via Hetzner Robot separately.",
	})
	d.SetId("")
	return diags
}
