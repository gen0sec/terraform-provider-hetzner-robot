package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataServerAddons() *schema.Resource {
	return &schema.Resource{
		Description: "Lists add-on products orderable for an existing server (GET /order/server_addon/{server-number}/product).",
		ReadContext: dataServerAddonsRead,
		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server number to list add-ons for",
			},
			"addons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":          {Type: schema.TypeString, Computed: true},
						"name":        {Type: schema.TypeString, Computed: true},
						"type":        {Type: schema.TypeString, Computed: true},
						"price_net":   {Type: schema.TypeString, Computed: true},
						"price_gross": {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataServerAddonsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber := d.Get("server_number").(int)

	addons, err := c.getServerAddonProducts(ctx, serverNumber)
	if err != nil {
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(addons))
	for _, a := range addons {
		result = append(result, map[string]interface{}{
			"id":          a.ID,
			"name":        a.Name,
			"type":        a.Type,
			"price_net":   a.PriceNet,
			"price_gross": a.PriceGross,
		})
	}

	if err := d.Set("addons", result); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(serverNumber))
	return nil
}
