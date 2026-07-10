package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#reverse-dns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type RdnsWrapper struct {
	Rdns Rdns `json:"rdns"`
}

type Rdns struct {
	IP  string `json:"ip"`
	PTR string `json:"ptr"`
}

func (c *HetznerRobotClient) getRdns(ctx context.Context, ip string) (*Rdns, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/rdns/%s", c.url, ip), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	wrapper := RdnsWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Rdns, nil
}

// setRdns creates or updates the PTR record for an IP (POST = create-or-update).
func (c *HetznerRobotClient) setRdns(ctx context.Context, ip string, ptr string) (*Rdns, error) {
	data := url.Values{}
	data.Set("ptr", ptr)

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/rdns/%s", c.url, ip), data, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	wrapper := RdnsWrapper{}
	if err = json.Unmarshal(bytes, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Rdns, nil
}

func (c *HetznerRobotClient) deleteRdns(ctx context.Context, ip string) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/rdns/%s", c.url, ip), nil, []int{http.StatusOK, http.StatusNotFound})
	return err
}
