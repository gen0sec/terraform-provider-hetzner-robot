package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataServerProducts() *schema.Resource {
	return &schema.Resource{
		Description: "Lists available Hetzner dedicated-server products for ordering (GET /order/server/product).",
		ReadContext: dataServerProductsRead,
		Schema: map[string]*schema.Schema{
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":          {Type: schema.TypeString, Computed: true},
						"name":        {Type: schema.TypeString, Computed: true},
						"description": {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"traffic":     {Type: schema.TypeString, Computed: true},
						"locations":   {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"prices": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location":      {Type: schema.TypeString, Computed: true},
									"monthly_net":   {Type: schema.TypeString, Computed: true},
									"monthly_gross": {Type: schema.TypeString, Computed: true},
									"setup_net":     {Type: schema.TypeString, Computed: true},
									"setup_gross":   {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataServerProductsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	products, err := c.getServerProducts(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		prices := make([]map[string]interface{}, 0, len(p.Prices))
		for _, pr := range p.Prices {
			prices = append(prices, map[string]interface{}{
				"location":      pr.Location,
				"monthly_net":   pr.MonthlyNet,
				"monthly_gross": pr.MonthlyGross,
				"setup_net":     pr.SetupNet,
				"setup_gross":   pr.SetupGross,
			})
		}
		result = append(result, map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"traffic":     p.Traffic,
			"locations":   p.Locations,
			"prices":      prices,
		})
	}

	if err := d.Set("products", result); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("server_products")
	return nil
}
