package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#failover

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type FailoverWrapper struct {
	Failover HetznerRobotFailover `json:"failover"`
}

type HetznerRobotFailover struct {
	IP             string `json:"ip"`
	Netmask        string `json:"netmask"`
	ServerIP       string `json:"server_ip"`
	ServerNumber   int    `json:"server_number"`
	ActiveServerIP string `json:"active_server_ip"`
}

func (c *HetznerRobotClient) getFailover(ctx context.Context, failoverIP string) (*HetznerRobotFailover, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/failover/%s", c.url, failoverIP), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := FailoverWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Failover, nil
}

// setFailover routes a failover IP to the given active server IP.
func (c *HetznerRobotClient) setFailover(ctx context.Context, failoverIP string, activeServerIP string) (*HetznerRobotFailover, error) {
	data := url.Values{}
	data.Set("active_server_ip", activeServerIP)

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/failover/%s", c.url, failoverIP), data, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := FailoverWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Failover, nil
}
