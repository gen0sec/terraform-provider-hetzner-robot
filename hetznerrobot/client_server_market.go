package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#server-ordering (server_market)

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type HetznerRobotMarketProduct struct {
	ID           string
	Name         string
	Description  []string
	Traffic      string
	Dist         []string
	Arch         []string
	Lang         []string
	CPU          string
	CPUBenchmark int
	MemorySize   int
	HDDSize      int
	HDDText      string
	HDDCount     int
	Datacenter   string
	NetworkSpeed string
	Price        string
	FixedPrice   bool
	NextReduce   int
}

func parseMarketProduct(p gjson.Result) HetznerRobotMarketProduct {
	return HetznerRobotMarketProduct{
		ID:           p.Get("id").String(),
		Name:         p.Get("name").String(),
		Description:  gjsonStrings(p.Get("description")),
		Traffic:      p.Get("traffic").String(),
		Dist:         gjsonStrings(p.Get("dist")),
		Arch:         gjsonStrings(p.Get("arch")),
		Lang:         gjsonStrings(p.Get("lang")),
		CPU:          p.Get("cpu").String(),
		CPUBenchmark: int(p.Get("cpu_benchmark").Int()),
		MemorySize:   int(p.Get("memory_size").Int()),
		HDDSize:      int(p.Get("hdd_size").Int()),
		HDDText:      p.Get("hdd_text").String(),
		HDDCount:     int(p.Get("hdd_count").Int()),
		Datacenter:   p.Get("datacenter").String(),
		NetworkSpeed: p.Get("network_speed").String(),
		Price:        p.Get("price").String(),
		FixedPrice:   p.Get("fixed_price").Bool(),
		NextReduce:   int(p.Get("next_reduce").Int()),
	}
}

func (c *HetznerRobotClient) getServerMarketProducts(ctx context.Context) ([]HetznerRobotMarketProduct, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server_market/product", c.url), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	products := []HetznerRobotMarketProduct{}
	gjson.ParseBytes(res).ForEach(func(_, item gjson.Result) bool {
		p := item.Get("product")
		if !p.Exists() {
			p = item
		}
		products = append(products, parseMarketProduct(p))
		return true
	})
	return products, nil
}

// createServerMarketOrder places an auction (server market) order. Market orders
// have no location — the listing is already tied to a datacenter.
func (c *HetznerRobotClient) createServerMarketOrder(ctx context.Context, req HetznerRobotServerOrderRequest) (*HetznerRobotServerTransaction, error) {
	data := url.Values{}
	data.Set("product_id", req.ProductID)
	if req.Dist != "" {
		data.Set("dist", req.Dist)
	}
	if req.Lang != "" {
		data.Set("lang", req.Lang)
	}
	for _, key := range req.AuthorizedKeys {
		data.Add("authorized_key[]", key)
	}
	if req.Password != "" {
		data.Set("password", req.Password)
	}
	for _, addon := range req.Addons {
		data.Add("addon[]", addon)
	}
	if req.Test {
		data.Set("test", "true")
	}

	res, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/order/server_market/transaction", c.url), data, []int{http.StatusOK, http.StatusCreated})
	if err != nil {
		return nil, err
	}
	return parseServerTransaction(string(res)), nil
}

func (c *HetznerRobotClient) getServerMarketOrder(ctx context.Context, id string) (*HetznerRobotServerTransaction, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/order/server_market/transaction/%s", c.url, id), nil, []int{http.StatusOK})
	if err != nil {
		return nil, err
	}
	return parseServerTransaction(string(res)), nil
}
