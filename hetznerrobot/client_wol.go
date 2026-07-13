package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#wake-on-lan

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// sendWOL sends a Wake-on-LAN packet to a server.
func (c *HetznerRobotClient) sendWOL(ctx context.Context, serverNumber int) error {
	_, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/wol/%d", c.url, serverNumber), url.Values{}, []int{http.StatusOK, http.StatusCreated})
	return err
}
