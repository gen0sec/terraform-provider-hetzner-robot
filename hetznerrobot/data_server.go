package hetznerrobot

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerRead,
		Schema: map[string]*schema.Schema{
			"server_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server name",
			},

			// read-only / computed
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data center",
			},
			"is_cancelled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Status of server cancellation",
			},
			"paid_until": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Paid until date",
			},
			"product": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server product name",
			},
			//"server_ip_addresses": {
			//	Type:        schema.TypeList,
			//	Computed:    true,
			//	Description: "Array of assigned single IP addresses",
			//	Elem: &schema.Schema{Type: schema.TypeString},
			//},
			//"server_ip_v4_addr": {
			//	Type:        schema.TypeString,
			//	Computed:    true,
			//	Description: "Server main IP address",
			//},
			//"server_ip_v6_net": {
			//	Type:        schema.TypeString,
			//	Computed:    true,
			//	Description: "Server main IPv6 net address",
			//},
			//"server_subnets": {
			//	Type:        schema.TypeList,
			//	Computed:    true,
			//	Description: "Array of assigned subnets",
			//	Elem: ,
			//},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server status (\"ready\" or \"in process\")",
			},
			"traffic": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Free traffic quota, 'unlimited' in case of unlimited traffic",
			},

			/*
			   reset (Boolean)	Flag of reset system availability
			   rescue (Boolean)	Flag of Rescue System availability
			   vnc (Boolean)	Flag of VNC installation availability
			   windows (Boolean)	Flag of Windows installation availability
			   plesk (Boolean)	Flag of Plesk installation availability
			   cpanel (Boolean)	Flag of cPanel installation availability
			   wol (Boolean)	Flag of Wake On Lan availability
			   hot_swap (Boolean)	Flag of Hot Swap availability

			   linked_storagebox (Integer)	Linked Storage Box ID
			*/
		},
	}
}

func dataSourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	serverID := d.Id()

	server, err := c.getServer(ctx, serverID)
	if err != nil {
		return diag.Errorf("Unable to find Server with IP %s:\n\t %q", serverID, err)
	}
	d.Set("datacenter", server.DataCenter)
	d.Set("is_cancelled", server.Cancelled)
	d.Set("paid_until", server.PaidUntil)
	d.Set("product", server.Product)
	//d.Set("server_ip_addresses", server.IPList)
	//d.Set("server_ip_v4_addr", server.ServerIPv4)
	//d.Set("server_ip_v6_net", server.ServerIPv6)
	d.Set("server_name", server.Name)
	//d.Set("server_subnets", server.SubnetList)
	d.Set("status", server.Status)
	d.Set("traffic", server.Traffic)
	d.SetId(serverID)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
