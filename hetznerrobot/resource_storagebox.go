package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceStorageBox manages the mutable settings of an existing Storage Box
// (name and service toggles). Storage Boxes themselves are ordered out of band;
// this resource adopts one by ID and reconciles its configuration.
func resourceStorageBox() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages the name and service toggles of an existing Storage Box.",
		CreateContext: resourceStorageBoxSet,
		ReadContext:   resourceStorageBoxRead,
		UpdateContext: resourceStorageBoxSet,
		DeleteContext: resourceIPNoOpDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"storagebox_id":         {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Storage Box ID"},
			"name":                  {Type: schema.TypeString, Required: true, Description: "Storage Box name"},
			"ssh":                   {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable SSH access"},
			"samba":                 {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable Samba access"},
			"webdav":                {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable WebDAV access"},
			"external_reachability": {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable external reachability"},
			"zfs":                   {Type: schema.TypeBool, Optional: true, Default: false, Description: "Enable ZFS snapshot directory visibility"},
			"login":                 {Type: schema.TypeString, Computed: true},
			"product":               {Type: schema.TypeString, Computed: true},
			"location":              {Type: schema.TypeString, Computed: true},
			"disk_quota":            {Type: schema.TypeInt, Computed: true},
			"disk_usage":            {Type: schema.TypeInt, Computed: true},
			"server":                {Type: schema.TypeString, Computed: true},
			"host_system":           {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceStorageBoxSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	b, err := c.updateStorageBox(ctx, id,
		d.Get("name").(string),
		d.Get("ssh").(bool),
		d.Get("samba").(bool),
		d.Get("webdav").(bool),
		d.Get("external_reachability").(bool),
		d.Get("zfs").(bool),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(id))
	setStorageBoxResourceAttrs(d, b)
	return nil
}

func resourceStorageBoxRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	b, err := c.getStorageBox(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("storagebox_id", b.ID)
	d.Set("name", b.Name)
	setStorageBoxResourceAttrs(d, b)
	return nil
}

func setStorageBoxResourceAttrs(d *schema.ResourceData, b *HetznerRobotStorageBox) {
	d.Set("ssh", b.SSH)
	d.Set("samba", b.Samba)
	d.Set("webdav", b.WebDAV)
	d.Set("external_reachability", b.ExternalReachability)
	d.Set("zfs", b.ZFS)
	d.Set("login", b.Login)
	d.Set("product", b.Product)
	d.Set("location", b.Location)
	d.Set("disk_quota", b.DiskQuota)
	d.Set("disk_usage", b.DiskUsage)
	d.Set("server", b.Server)
	d.Set("host_system", b.HostSystem)
}
