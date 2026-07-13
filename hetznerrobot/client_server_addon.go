package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#server-ordering (server_addon)

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tidwall/gjson"
)

type HetznerRobotAddonProduct struct {
	ID         string
	Name       string
	Type       string
	PriceNet   string
	PriceGross string
}

func (c *HetznerRobotClient) getServerAddonProducts(ctx context.Context, serverNumber int) ([]HetznerRobotAddonProduct, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server_addon/%d/product", c.url, serverNumber), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	addons := []HetznerRobotAddonProduct{}
	gjson.ParseBytes(res).ForEach(func(_, item gjson.Result) bool {
		p := item.Get("product")
		if !p.Exists() {
			p = item
		}
		addons = append(addons, HetznerRobotAddonProduct{
			ID:         p.Get("id").String(),
			Name:       p.Get("name").String(),
			Type:       p.Get("type").String(),
			PriceNet:   p.Get("price.price.net").String(),
			PriceGross: p.Get("price.price.gross").String(),
		})
		return true
	})
	return addons, nil
}

// createServerAddonOrder orders an add-on for an existing server. When test is
// true the order is validated only (no charge). reason is an optional
// justification required by some add-ons (e.g. an additional/primary IPv4).
func (c *HetznerRobotClient) createServerAddonOrder(ctx context.Context, serverNumber int, productID, reason string, test bool) (*HetznerRobotServerTransaction, error) {
	data := url.Values{}
	data.Set("server_number", strconv.Itoa(serverNumber))
	data.Set("product_id", productID)
	if reason != "" {
		data.Set("reason", reason)
	}
	if test {
		data.Set("test", "true")
	}

	res, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/order/server_addon/transaction", c.url), data, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	return parseServerTransaction(string(res)), nil
}

func (c *HetznerRobotClient) getServerAddonOrder(ctx context.Context, id string) (*HetznerRobotServerTransaction, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server_addon/transaction/%s", c.url, id), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	return parseServerTransaction(string(res)), nil
}
