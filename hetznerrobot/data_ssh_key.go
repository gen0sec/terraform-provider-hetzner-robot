package hetznerrobot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSshKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSshKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Key name. Provide either name or fingerprint to look up a key.",
				ExactlyOneOf: []string{"name", "fingerprint"},
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key data in OpenSSH or SSH2 format",
			},
			"fingerprint": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Key fingerprint. Provide either name or fingerprint to look up a key.",
				ExactlyOneOf: []string{"name", "fingerprint"},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key algorithm type",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Key size in bits",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date",
			},
		},
	}
}

func dataSourceSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	var key *SshKey
	var err error
	if fingerprint := d.Get("fingerprint").(string); fingerprint != "" {
		key, err = c.getSshKey(ctx, fingerprint)
	} else {
		key, err = c.getSshKeyByName(ctx, d.Get("name").(string))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", key.Name)
	d.Set("data", key.Data)
	d.Set("fingerprint", key.Fingerprint)
	d.Set("type", key.Type)
	d.Set("size", key.Size)
	d.Set("created_at", key.CreatedAt)

	d.SetId(key.Fingerprint)

	return diag.Diagnostics{}
}
