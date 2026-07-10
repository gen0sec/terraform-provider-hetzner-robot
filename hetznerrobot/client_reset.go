package hetznerrobot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// resetServer triggers a reset of the given server via the Robot API. Valid
// types are "hw" (hardware reset), "sw" (Ctrl+Alt+Del), "power" (power cycle)
// and "man" (manual, performed by support).
func (c *HetznerRobotClient) resetServer(ctx context.Context, serverID int, resetType string) error {
	data := url.Values{}
	data.Set("type", resetType)
	_, err := c.makeAPICall(
		ctx,
		"POST",
		fmt.Sprintf("%s/reset/%d", c.url, serverID),
		data,
		[]int{http.StatusOK, http.StatusAccepted},
	)
	return err
}
