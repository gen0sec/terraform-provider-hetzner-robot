package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#server-cancellation

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type HetznerRobotCancellation struct {
	ServerNumber             int
	ServerIP                 string
	CancellationDate         string
	Cancelled                bool
	EarliestCancellationDate string
}

func parseCancellation(jsonStr string) *HetznerRobotCancellation {
	c := gjson.Get(jsonStr, "cancellation")
	if !c.Exists() {
		c = gjson.Parse(jsonStr)
	}
	return &HetznerRobotCancellation{
		ServerNumber:             int(c.Get("server_number").Int()),
		ServerIP:                 c.Get("server_ip").String(),
		CancellationDate:         c.Get("cancellation_date").String(),
		Cancelled:                c.Get("cancelled").Bool(),
		EarliestCancellationDate: c.Get("earliest_cancellation_date").String(),
	}
}

func (c *HetznerRobotClient) getCancellation(ctx context.Context, serverNumber int) (*HetznerRobotCancellation, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/server/%d/cancellation", c.url, serverNumber), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	return parseCancellation(string(bytes)), nil
}

// cancelServer schedules cancellation of a server. cancellationDate is "now" or
// a date (YYYY-MM-DD); reason is optional.
func (c *HetznerRobotClient) cancelServer(ctx context.Context, serverNumber int, cancellationDate string, reason string) (*HetznerRobotCancellation, error) {
	data := url.Values{}
	data.Set("cancellation_date", cancellationDate)
	if reason != "" {
		data.Set("cancellation_reason", reason)
	}

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/server/%d/cancellation", c.url, serverNumber), data, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	return parseCancellation(string(bytes)), nil
}

// revokeCancellation withdraws a pending cancellation.
func (c *HetznerRobotClient) revokeCancellation(ctx context.Context, serverNumber int) error {
	_, err := c.makeAPICall(ctx, "DELETE", fmt.Sprintf("%s/server/%d/cancellation", c.url, serverNumber), nil, []int{http.StatusOK, http.StatusNotFound})
	return err
}
