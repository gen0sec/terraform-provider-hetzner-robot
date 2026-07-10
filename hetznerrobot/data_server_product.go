package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataServerProduct() *schema.Resource {
	return &schema.Resource{
		Description: "Details and orderable options for a single Hetzner server product (GET /order/server/product/{id}). Use the dist/lang/location lists to fill a hetzner-robot_server_order.",
		ReadContext: dataServerProductRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product ID (e.g. AX42-1).",
			},
			"name":        {Type: schema.TypeString, Computed: true, Description: "Product name"},
			"description": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Specification lines"},
			"traffic":     {Type: schema.TypeString, Computed: true, Description: "Traffic allowance"},
			"locations":   {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Available locations"},
			"dist":        {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Orderable distributions"},
			"lang":        {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Orderable languages"},
			"arch":        {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}, Description: "Orderable architectures"},
		},
	}
}

func dataServerProductRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	p, err := c.getServerProduct(ctx, d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", p.Name)
	d.Set("description", p.Description)
	d.Set("traffic", p.Traffic)
	d.Set("locations", p.Locations)
	d.Set("dist", p.Dist)
	d.Set("lang", p.Lang)
	d.Set("arch", p.Arch)
	d.SetId(p.ID)
	return nil
}
