package hetznerrobot

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStorageBoxSnapshot() *schema.Resource {
	return &schema.Resource{
		Description:   "Creates a snapshot of a Storage Box.",
		CreateContext: resourceStorageBoxSnapshotCreate,
		ReadContext:   resourceStorageBoxSnapshotRead,
		DeleteContext: resourceStorageBoxSnapshotDelete,
		Schema: map[string]*schema.Schema{
			"storagebox_id": {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Storage Box ID"},
			"name":          {Type: schema.TypeString, Computed: true, Description: "Snapshot name assigned by Hetzner"},
			"timestamp":     {Type: schema.TypeString, Computed: true},
			"size":          {Type: schema.TypeInt, Computed: true},
			"automatic":     {Type: schema.TypeBool, Computed: true},
		},
	}
}

func resourceStorageBoxSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	snap, err := c.createStorageBoxSnapshot(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d/%s", id, snap.Name))
	d.Set("name", snap.Name)
	d.Set("timestamp", snap.Timestamp)
	d.Set("size", snap.Size)
	d.Set("automatic", snap.Automatic)
	return nil
}

func resourceStorageBoxSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	name := d.Get("name").(string)
	snaps, err := c.getStorageBoxSnapshots(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, s := range snaps {
		if s.Name == name {
			d.Set("timestamp", s.Timestamp)
			d.Set("size", s.Size)
			d.Set("automatic", s.Automatic)
			return nil
		}
	}
	// snapshot no longer exists
	d.SetId("")
	return nil
}

func resourceStorageBoxSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	if err := c.deleteStorageBoxSnapshot(ctx, id, d.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
