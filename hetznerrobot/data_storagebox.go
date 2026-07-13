package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func storageBoxComputedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"storagebox_id":         {Type: schema.TypeInt, Computed: true},
		"login":                 {Type: schema.TypeString, Computed: true},
		"name":                  {Type: schema.TypeString, Computed: true},
		"product":               {Type: schema.TypeString, Computed: true},
		"cancelled":             {Type: schema.TypeBool, Computed: true},
		"locked":                {Type: schema.TypeBool, Computed: true},
		"location":              {Type: schema.TypeString, Computed: true},
		"linked_server":         {Type: schema.TypeInt, Computed: true},
		"paid_until":            {Type: schema.TypeString, Computed: true},
		"disk_quota":            {Type: schema.TypeInt, Computed: true},
		"disk_usage":            {Type: schema.TypeInt, Computed: true},
		"webdav":                {Type: schema.TypeBool, Computed: true},
		"samba":                 {Type: schema.TypeBool, Computed: true},
		"ssh":                   {Type: schema.TypeBool, Computed: true},
		"external_reachability": {Type: schema.TypeBool, Computed: true},
		"zfs":                   {Type: schema.TypeBool, Computed: true},
		"server":                {Type: schema.TypeString, Computed: true},
		"host_system":           {Type: schema.TypeString, Computed: true},
	}
}

func setStorageBoxData(d *schema.ResourceData, b *HetznerRobotStorageBox) {
	d.Set("storagebox_id", b.ID)
	d.Set("login", b.Login)
	d.Set("name", b.Name)
	d.Set("product", b.Product)
	d.Set("cancelled", b.Cancelled)
	d.Set("locked", b.Locked)
	d.Set("location", b.Location)
	d.Set("linked_server", b.LinkedServer)
	d.Set("paid_until", b.PaidUntil)
	d.Set("disk_quota", b.DiskQuota)
	d.Set("disk_usage", b.DiskUsage)
	d.Set("webdav", b.WebDAV)
	d.Set("samba", b.Samba)
	d.Set("ssh", b.SSH)
	d.Set("external_reachability", b.ExternalReachability)
	d.Set("zfs", b.ZFS)
	d.Set("server", b.Server)
	d.Set("host_system", b.HostSystem)
}

func storageBoxMap(b *HetznerRobotStorageBox) map[string]interface{} {
	return map[string]interface{}{
		"storagebox_id":         b.ID,
		"login":                 b.Login,
		"name":                  b.Name,
		"product":               b.Product,
		"cancelled":             b.Cancelled,
		"locked":                b.Locked,
		"location":              b.Location,
		"linked_server":         b.LinkedServer,
		"paid_until":            b.PaidUntil,
		"disk_quota":            b.DiskQuota,
		"disk_usage":            b.DiskUsage,
		"webdav":                b.WebDAV,
		"samba":                 b.Samba,
		"ssh":                   b.SSH,
		"external_reachability": b.ExternalReachability,
		"zfs":                   b.ZFS,
		"server":                b.Server,
		"host_system":           b.HostSystem,
	}
}

func dataStorageBox() *schema.Resource {
	s := storageBoxComputedSchema()
	s["storagebox_id"] = &schema.Schema{Type: schema.TypeInt, Required: true, Description: "Storage Box ID"}
	return &schema.Resource{
		Description: "Reads a single Storage Box by ID.",
		ReadContext: dataStorageBoxRead,
		Schema:      s,
	}
}

func dataStorageBoxRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	b, err := c.getStorageBox(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	setStorageBoxData(d, b)
	d.SetId(strconv.Itoa(b.ID))
	return nil
}
