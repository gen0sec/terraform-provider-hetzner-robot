package hetznerrobot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HetznerRobotVSwitchServer struct {
	ServerIP      string `json:"server_ip"`
	ServerIPv6Net string `json:"server_ipv6_net"`
	ServerNumber  int    `json:"server_number"`
	Status        string `json:"status"`
}

type HetznerRobotVSwitchSubnet struct {
	IP      string `json:"ip"`
	Mask    int    `json:"mask"`
	Gateway string `json:"gateway"`
}

type HetznerRobotVSwitchCloudNetwork struct {
	ID      int    `json:"id"`
	IP      string `json:"ip"`
	Mask    int    `json:"mask"`
	Gateway string `json:"gateway"`
}

type HetznerRobotVSwitch struct {
	ID           int                               `json:"id"`
	Name         string                            `json:"name"`
	Vlan         int                               `json:"vlan"`
	Cancelled    bool                              `json:"cancelled"`
	Server       []HetznerRobotVSwitchServer       `json:"server"`
	Subnet       []HetznerRobotVSwitchSubnet       `json:"subnet"`
	CloudNetwork []HetznerRobotVSwitchCloudNetwork `json:"cloud_network"`
}

func (c *HetznerRobotClient) getVSwitch(ctx context.Context, id string) (*HetznerRobotVSwitch, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/vswitch/%s", c.url, id), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}

	vSwitch := HetznerRobotVSwitch{}
	if err = json.Unmarshal(res, &vSwitch); err != nil {
		return nil, err
	}
	return &vSwitch, nil
}

func (c *HetznerRobotClient) createVSwitch(ctx context.Context, name string, vlan int) (*HetznerRobotVSwitch, error) {
	data := url.Values{}
	data.Set("vlan", strconv.Itoa(vlan))
	data.Set("name", name)
	res, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/vswitch", c.url), strings.NewReader(data.Encode()), http.StatusOK)
	if err != nil {
		return nil, err
	}

	vSwitch := HetznerRobotVSwitch{}
	if err = json.Unmarshal(res, &vSwitch); err != nil {
		return nil, err
	}
	return &vSwitch, nil
}

func (c *HetznerRobotClient) updateVSwitch(ctx context.Context, id string, name string, vlan int) error {
	data := url.Values{}
	data.Set("vlan", strconv.Itoa(vlan))
	data.Set("name", name)
	_, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/vswitch/%s", c.url, id), strings.NewReader(data.Encode()), http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}

func (c *HetznerRobotClient) deleteVSwitch(ctx context.Context, id string) error {
	now := time.Now()
	data := url.Values{}
	data.Set("cancellation_date", now.Format("2006-01-02"))
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/vswitch/%s", c.url, id), strings.NewReader(data.Encode()), http.StatusOK)
	if err != nil {
		return err
	}
	return nil
}
