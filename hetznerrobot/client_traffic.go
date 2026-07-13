package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#traffic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type HetznerRobotTraffic struct {
	In  string
	Out string
	Sum string
}

// getTraffic queries traffic for a single IP over a period. type is day/month/year.
func (c *HetznerRobotClient) getTraffic(ctx context.Context, ip, trafficType, from, to string) (*HetznerRobotTraffic, error) {
	data := url.Values{}
	data.Set("type", trafficType)
	data.Set("from", from)
	data.Set("to", to)
	data.Add("ip[]", ip)

	res, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/traffic", c.url), data, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}

	// data is keyed by IP (whose dots would break a gjson path), so take the
	// first/only entry.
	t := &HetznerRobotTraffic{}
	gjson.GetBytes(res, "traffic.data").ForEach(func(_, val gjson.Result) bool {
		t.In = val.Get("in").String()
		t.Out = val.Get("out").String()
		t.Sum = val.Get("sum").String()
		return false
	})
	return t, nil
}
