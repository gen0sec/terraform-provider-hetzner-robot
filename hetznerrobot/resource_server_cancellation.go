package hetznerrobot

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServerCancellation() *schema.Resource {
	return &schema.Resource{
		Description: "Schedules cancellation of a dedicated server. Creating this resource cancels the server; " +
			"destroying it revokes the cancellation. This is the counterpart to hetzner-robot_server_order.",
		CreateContext: resourceServerCancellationSet,
		ReadContext:   resourceServerCancellationRead,
		UpdateContext: resourceServerCancellationSet,
		DeleteContext: resourceServerCancellationDelete,
		Schema: map[string]*schema.Schema{
			"server_number": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Server number to cancel",
			},
			"cancellation_date": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cancellation date: 'now' or a date (YYYY-MM-DD)",
			},
			"cancellation_reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional cancellation reason",
			},

			"cancelled":                  {Type: schema.TypeBool, Computed: true, Description: "Whether the server is cancelled"},
			"earliest_cancellation_date": {Type: schema.TypeString, Computed: true, Description: "Earliest possible cancellation date"},
		},
	}
}

func resourceServerCancellationSet(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber := d.Get("server_number").(int)

	cancellation, err := c.cancelServer(ctx, serverNumber, d.Get("cancellation_date").(string), d.Get("cancellation_reason").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(serverNumber))
	d.Set("cancelled", cancellation.Cancelled)
	d.Set("earliest_cancellation_date", cancellation.EarliestCancellationDate)
	return nil
}

func resourceServerCancellationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	cancellation, err := c.getCancellation(ctx, serverNumber)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("server_number", cancellation.ServerNumber)
	d.Set("cancelled", cancellation.Cancelled)
	d.Set("earliest_cancellation_date", cancellation.EarliestCancellationDate)
	if cancellation.CancellationDate != "" {
		d.Set("cancellation_date", cancellation.CancellationDate)
	}
	return nil
}

func resourceServerCancellationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)
	serverNumber, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := c.revokeCancellation(ctx, serverNumber); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
