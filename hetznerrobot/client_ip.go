package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#ip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type IPWrapper struct {
	IP HetznerRobotIP `json:"ip"`
}

type HetznerRobotIP struct {
	IP              string `json:"ip"`
	ServerIP        string `json:"server_ip"`
	ServerNumber    int    `json:"server_number"`
	Locked          bool   `json:"locked"`
	SeparateMAC     string `json:"separate_mac"`
	TrafficWarnings bool   `json:"traffic_warnings"`
	TrafficHourly   int    `json:"traffic_hourly"`
	TrafficDaily    int    `json:"traffic_daily"`
	TrafficMonthly  int    `json:"traffic_monthly"`
}

func (c *HetznerRobotClient) getIP(ctx context.Context, ip string) (*HetznerRobotIP, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/ip/%s", c.url, ip), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := IPWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.IP, nil
}

// setIP updates an IP's traffic-warning configuration.
func (c *HetznerRobotClient) setIP(ctx context.Context, ip string, warnings bool, hourly, daily, monthly int) (*HetznerRobotIP, error) {
	data := url.Values{}
	data.Set("traffic_warnings", strconv.FormatBool(warnings))
	data.Set("traffic_hourly", strconv.Itoa(hourly))
	data.Set("traffic_daily", strconv.Itoa(daily))
	data.Set("traffic_monthly", strconv.Itoa(monthly))

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/ip/%s", c.url, ip), data, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := IPWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.IP, nil
}

// --- virtual MAC (vMAC) for a single IP ---

type MacWrapper struct {
	Mac HetznerRobotMac `json:"mac"`
}

type HetznerRobotMac struct {
	IP       string `json:"ip"`
	Mac      string `json:"mac"`
	Possible bool   `json:"possible"`
}

func (c *HetznerRobotClient) getIPMac(ctx context.Context, ip string) (*HetznerRobotMac, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/ip/%s/mac", c.url, ip), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := MacWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Mac, nil
}

func (c *HetznerRobotClient) createIPMac(ctx context.Context, ip string) (*HetznerRobotMac, error) {
	bytes, err := c.makeAPICall(ctx, "PUT", fmt.Sprintf("%s/ip/%s/mac", c.url, ip), url.Values{}, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	wrapper := MacWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Mac, nil
}

func (c *HetznerRobotClient) deleteIPMac(ctx context.Context, ip string) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/ip/%s/mac", c.url, ip), nil, []int{http.StatusOK, http.StatusNotFound})
	return err
}
