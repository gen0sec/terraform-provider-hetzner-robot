package hetznerrobot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HetznerRobotServerResponse struct {
	Server HetznerRobotServer `json:"server"`
}

type HetznerRobotServer struct {
	IP         string `json:"server_ip"`
	Number     int    `json:"server_number"`
	Name       string `json:"server_name"`
	Product    string `json:"product"`
	DataCenter string `json:"dc"`
	Traffic    string `json:"traffic"`
	Status     string `json:"status"`
	Cancelled  bool   `json:"cancelled"`
	PaidUntil  string `json:"paid_until"`
}

type HetznerRobotServerRenameRequestBody struct {
	Name string `json:"server_name"`
}

func (c *HetznerRobotClient) getServer(ctx context.Context, ip string) (*HetznerRobotServer, error) {

	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/server/%s", c.url, ip), nil, http.StatusOK)
	if err != nil {
		return nil, err
	}

	serverResponse := HetznerRobotServerResponse{}
	if err = json.Unmarshal(res, &serverResponse); err != nil {
		return nil, err
	}
	return &serverResponse.Server, nil
}

func (c *HetznerRobotClient) setServerName(ctx context.Context, ip string, name string) (*HetznerRobotServer, error) {
	body, _ := json.Marshal(&HetznerRobotServerRenameRequestBody{Name: name})
	res, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/server/%s", c.url, ip), bytes.NewReader(body), http.StatusOK)
	if err != nil {
		return nil, err
	}

	serverResponse := HetznerRobotServerResponse{}
	if err = json.Unmarshal(res, &serverResponse); err != nil {
		return nil, err
	}
	return &serverResponse.Server, nil
}
