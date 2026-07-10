package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataServerMarketProducts() *schema.Resource {
	return &schema.Resource{
		Description: "Lists available Hetzner server market (auction) products (GET /order/server_market/product).",
		ReadContext: dataServerMarketProductsRead,
		Schema: map[string]*schema.Schema{
			"products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            {Type: schema.TypeString, Computed: true},
						"name":          {Type: schema.TypeString, Computed: true},
						"description":   {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"traffic":       {Type: schema.TypeString, Computed: true},
						"dist":          {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"lang":          {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"arch":          {Type: schema.TypeList, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
						"cpu":           {Type: schema.TypeString, Computed: true},
						"cpu_benchmark": {Type: schema.TypeInt, Computed: true},
						"memory_size":   {Type: schema.TypeInt, Computed: true},
						"hdd_size":      {Type: schema.TypeInt, Computed: true},
						"hdd_text":      {Type: schema.TypeString, Computed: true},
						"hdd_count":     {Type: schema.TypeInt, Computed: true},
						"datacenter":    {Type: schema.TypeString, Computed: true},
						"network_speed": {Type: schema.TypeString, Computed: true},
						"price":         {Type: schema.TypeString, Computed: true},
						"fixed_price":   {Type: schema.TypeBool, Computed: true},
						"next_reduce":   {Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataServerMarketProductsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	products, err := c.getServerMarketProducts(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		result = append(result, map[string]interface{}{
			"id":            p.ID,
			"name":          p.Name,
			"description":   p.Description,
			"traffic":       p.Traffic,
			"dist":          p.Dist,
			"lang":          p.Lang,
			"arch":          p.Arch,
			"cpu":           p.CPU,
			"cpu_benchmark": p.CPUBenchmark,
			"memory_size":   p.MemorySize,
			"hdd_size":      p.HDDSize,
			"hdd_text":      p.HDDText,
			"hdd_count":     p.HDDCount,
			"datacenter":    p.Datacenter,
			"network_speed": p.NetworkSpeed,
			"price":         p.Price,
			"fixed_price":   p.FixedPrice,
			"next_reduce":   p.NextReduce,
		})
	}

	if err := d.Set("products", result); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("server_market_products")
	return nil
}
