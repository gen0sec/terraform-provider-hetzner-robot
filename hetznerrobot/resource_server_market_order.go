package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerMarketOrder() *schema.Resource {
	return &schema.Resource{
		Description: "Orders a Hetzner server market (auction) server (POST /order/server_market/transaction). " +
			"By default test = true, which only validates the order without provisioning or charging. Set " +
			"test = false to place a REAL, billable order. Destroying this resource does NOT cancel the server.",
		CreateContext: resourceServerMarketOrderCreate,
		ReadContext:   resourceServerMarketOrderRead,
		DeleteContext: resourceServerOrderDelete, // same no-op-with-warning as the regular order
		Schema: map[string]*schema.Schema{
			"product_id":      {Type: schema.TypeString, Required: true, ForceNew: true, Description: "Auction product ID (see the hetzner-robot_server_market_products data source)"},
			"dist":            {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Distribution/OS to preinstall"},
			"lang":            {Type: schema.TypeString, Optional: true, ForceNew: true, Description: "Installation language"},
			"authorized_keys": {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "SSH key fingerprints authorized on the server"},
			"password":        {Type: schema.TypeString, Optional: true, ForceNew: true, Sensitive: true, Description: "Root password (alternative to authorized_keys)"},
			"addons":          {Type: schema.TypeList, Optional: true, ForceNew: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Product add-on IDs"},
			"primary_ipv4": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Convenience flag: when true, adds the \"primary_ipv4\" add-on so the auction server is ordered with a public IPv4 (auction servers are IPv6-only by default). Equivalent to including \"primary_ipv4\" in addons.",
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

func resourceServerMarketOrderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(HetznerRobotClient)

	req := HetznerRobotServerOrderRequest{
		ProductID:      d.Get("product_id").(string),
		Dist:           d.Get("dist").(string),
		Lang:           d.Get("lang").(string),
		AuthorizedKeys: serverOrderStringList(d, "authorized_keys"),
		Password:       d.Get("password").(string),
		Addons:         serverOrderAddons(d),
		Test:           d.Get("test").(bool),
	}

	tx, err := c.createServerMarketOrder(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if tx.ID != "" {
		d.SetId(tx.ID)
	} else {
		d.SetId("market-order-" + req.ProductID)
	}
	d.Set("transaction_id", tx.ID)
	d.Set("status", tx.Status)
	d.Set("server_number", tx.ServerNumber)
	d.Set("server_ip", tx.ServerIP)

	if req.Test {
		diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: "test order only (no server provisioned)", Detail: "test = true validated the order without provisioning or charging. Set test = false to place a real order."})
	} else {
		diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: "real, billable auction order placed", Detail: "A billable server-market order was placed. Destroying this resource does NOT cancel the server."})
	}
	return diags
}

func resourceServerMarketOrderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	tx, err := c.getServerMarketOrder(ctx, d.Id())
	if err != nil {
		return nil // test/expired transactions may not be retrievable
	}
	d.Set("transaction_id", tx.ID)
	d.Set("status", tx.Status)
	d.Set("server_number", tx.ServerNumber)
	d.Set("server_ip", tx.ServerIP)
	return nil
}
