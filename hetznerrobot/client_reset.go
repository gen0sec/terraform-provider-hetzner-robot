package hetznerrobot

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// resetServer triggers a reset of the given server via the Robot API. Valid
// types are "hw" (hardware reset), "sw" (Ctrl+Alt+Del), "power" (power cycle)
// and "man" (manual, performed by support).
func (c *HetznerRobotClient) resetServer(serverID int, resetType string) error {
	formParams := url.Values{}
	formParams.Set("type", resetType)
	_, err := c.makeAPICall(
		"POST",
		fmt.Sprintf("%s/reset/%d", c.url, serverID),
		strings.NewReader(formParams.Encode()),
		http.StatusOK,
	)
	return err
}
