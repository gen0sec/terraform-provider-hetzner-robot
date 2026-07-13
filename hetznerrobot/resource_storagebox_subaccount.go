package hetznerrobot

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStorageBoxSubaccount() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a subaccount of a Storage Box.",
		CreateContext: resourceStorageBoxSubaccountCreate,
		ReadContext:   resourceStorageBoxSubaccountRead,
		UpdateContext: resourceStorageBoxSubaccountUpdate,
		DeleteContext: resourceStorageBoxSubaccountDelete,
		Schema: map[string]*schema.Schema{
			"storagebox_id":         {Type: schema.TypeInt, Required: true, ForceNew: true, Description: "Storage Box ID"},
			"homedirectory":         {Type: schema.TypeString, Required: true, Description: "Home directory of the subaccount"},
			"samba":                 {Type: schema.TypeBool, Optional: true, Default: false},
			"ssh":                   {Type: schema.TypeBool, Optional: true, Default: false},
			"external_reachability": {Type: schema.TypeBool, Optional: true, Default: false},
			"webdav":                {Type: schema.TypeBool, Optional: true, Default: false},
			"readonly":              {Type: schema.TypeBool, Optional: true, Default: false},
			"comment":               {Type: schema.TypeString, Optional: true, Default: ""},
			"username":              {Type: schema.TypeString, Computed: true, Description: "Username assigned by Hetzner"},
			"password":              {Type: schema.TypeString, Computed: true, Sensitive: true, Description: "Password generated on creation"},
			"server":                {Type: schema.TypeString, Computed: true},
		},
	}
}

func resourceStorageBoxSubaccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	sub, err := c.createStorageBoxSubaccount(ctx, id,
		d.Get("homedirectory").(string),
		d.Get("samba").(bool),
		d.Get("ssh").(bool),
		d.Get("external_reachability").(bool),
		d.Get("webdav").(bool),
		d.Get("readonly").(bool),
		d.Get("comment").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d/%s", id, sub.Username))
	d.Set("username", sub.Username)
	d.Set("password", sub.Password)
	d.Set("server", sub.Server)
	return nil
}

func resourceStorageBoxSubaccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	username := d.Get("username").(string)
	subs, err := c.getStorageBoxSubaccounts(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, s := range subs {
		if s.Username == username {
			d.Set("homedirectory", s.HomeDirectory)
			d.Set("samba", s.Samba)
			d.Set("ssh", s.SSH)
			d.Set("external_reachability", s.ExternalReachability)
			d.Set("webdav", s.WebDAV)
			d.Set("readonly", s.Readonly)
			d.Set("comment", s.Comment)
			d.Set("server", s.Server)
			return nil
		}
	}
	d.SetId("")
	return nil
}

func resourceStorageBoxSubaccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	if err := c.updateStorageBoxSubaccount(ctx, id, d.Get("username").(string),
		d.Get("homedirectory").(string),
		d.Get("samba").(bool),
		d.Get("ssh").(bool),
		d.Get("external_reachability").(bool),
		d.Get("webdav").(bool),
		d.Get("readonly").(bool),
		d.Get("comment").(string),
	); err != nil {
		return diag.FromErr(err)
	}
	return resourceStorageBoxSubaccountRead(ctx, d, meta)
}

func resourceStorageBoxSubaccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	id := d.Get("storagebox_id").(int)
	if err := c.deleteStorageBoxSubaccount(ctx, id, d.Get("username").(string)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
