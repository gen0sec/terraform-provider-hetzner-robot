package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#subnet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SubnetWrapper struct {
	Subnet HetznerRobotSubnet `json:"subnet"`
}

type HetznerRobotSubnet struct {
	IP              string `json:"ip"`
	Mask            int    `json:"mask"`
	Gateway         string `json:"gateway"`
	ServerIP        string `json:"server_ip"`
	ServerNumber    int    `json:"server_number"`
	Failover        bool   `json:"failover"`
	Locked          bool   `json:"locked"`
	TrafficWarnings bool   `json:"traffic_warnings"`
	TrafficHourly   int    `json:"traffic_hourly"`
	TrafficDaily    int    `json:"traffic_daily"`
	TrafficMonthly  int    `json:"traffic_monthly"`
}

func (c *HetznerRobotClient) getSubnet(ctx context.Context, netIP string) (*HetznerRobotSubnet, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/subnet/%s", c.url, netIP), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := SubnetWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Subnet, nil
}

// setSubnet updates a subnet's traffic-warning configuration.
func (c *HetznerRobotClient) setSubnet(ctx context.Context, netIP string, warnings bool, hourly, daily, monthly int) (*HetznerRobotSubnet, error) {
	data := url.Values{}
	data.Set("traffic_warnings", strconv.FormatBool(warnings))
	data.Set("traffic_hourly", strconv.Itoa(hourly))
	data.Set("traffic_daily", strconv.Itoa(daily))
	data.Set("traffic_monthly", strconv.Itoa(monthly))

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/subnet/%s", c.url, netIP), data, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := SubnetWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Subnet, nil
}

// --- virtual MAC (vMAC) for a subnet ---

func (c *HetznerRobotClient) getSubnetMac(ctx context.Context, netIP string) (*HetznerRobotMac, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/subnet/%s/mac", c.url, netIP), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := MacWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Mac, nil
}

func (c *HetznerRobotClient) createSubnetMac(ctx context.Context, netIP string) (*HetznerRobotMac, error) {
	bytes, err := c.makeAPICall(ctx, "PUT", fmt.Sprintf("%s/subnet/%s/mac", c.url, netIP), url.Values{}, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	wrapper := MacWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Mac, nil
}

func (c *HetznerRobotClient) deleteSubnetMac(ctx context.Context, netIP string) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/subnet/%s/mac", c.url, netIP), nil, []int{http.StatusOK, http.StatusNotFound})
	return err
}
