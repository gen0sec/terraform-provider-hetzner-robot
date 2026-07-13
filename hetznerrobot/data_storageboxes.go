package hetznerrobot

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataStorageBoxes() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all Storage Boxes in the account.",
		ReadContext: dataStorageBoxesRead,
		Schema: map[string]*schema.Schema{
			"storageboxes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Resource{Schema: storageBoxComputedSchema()},
			},
		},
	}
}

func dataStorageBoxesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	boxes, err := c.getStorageBoxes(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	items := make([]map[string]interface{}, len(boxes))
	for i := range boxes {
		items[i] = storageBoxMap(&boxes[i])
	}
	if err := d.Set("storageboxes", items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("storageboxes")
	return nil
}
